package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/anthony-dong/aliyun-oss-cli/config"
	"github.com/anthony-dong/aliyun-oss-cli/logger"
	"github.com/anthony-dong/aliyun-oss-cli/util"
	"github.com/satori/go.uuid"
	"os"
	"strings"
	"sync"
)

type Client struct {
	sync.RWMutex
	Client         *oss.Client
	Config         *config.Config
	MakeUrl        MakeUrlFunc
	bucketName     string
	selectedBucket *oss.Bucket
}

func (this *Client) GetSelectedBucket() *oss.Bucket {
	this.RLock()
	defer this.RUnlock()
	return this.selectedBucket
}

type MakeUrlFunc func(fileName string) string // 创建url
var (
	cli = new(Client)
)

func GetClient() *Client {
	return cli
}

const (
	ImagePrefix string = "image/"
	FilePrefix  string = "file/"
)

func init() {
	cli.Config = config.GetOssConfig()
	var (
		urlPrefix string
		styleName string
	)
	if config.GetOssConfig().FastEndpoint != "" {
		urlPrefix = "https://" + cli.Config.BucketName + "." + config.GetOssConfig().FastEndpoint + "/"
	} else {
		urlPrefix = "https://" + cli.Config.BucketName + "." + config.GetOssConfig().Endpoint + "/"
	}
	if config.GetOssConfig().StyleName != "" {
		styleName = "?" + strings.TrimSpace(cli.Config.StyleName)
	}
	cli.MakeUrl = func(fileName string) string {
		return fmt.Sprintf("%s%s%s", urlPrefix, fileName, styleName)
	}
	client, err := oss.New(cli.Config.Endpoint, cli.Config.AccessKeyId, cli.Config.AccessKeySecret)
	if err != nil {
		logger.FatalF("conn oss server err,err:%s", err)
	}
	cli.Client = client
}

func (this *Client) GetAndNewBucket(bucketName string) *oss.Bucket {
	if err := oss.CheckBucketName(bucketName); err != nil {
		err := this.Client.CreateBucket(bucketName)
		if err != nil {
			logger.FatalF("create bucket err,err:%s", err)
		}
	}
	bucket, err := this.Client.Bucket(bucketName)
	if err != nil {
		logger.FatalF("open bucket err,err:%s", err)
	}
	return bucket
}

func (this *Client) ListBucket() []string {
	list := make([]string, 0)
	// 列举存储空间。
	marker := ""
	for {
		lsRes, err := this.Client.ListBuckets(oss.Marker(marker))
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		// 默认情况下一次返回100条记录。
		for _, bucket := range lsRes.Buckets {
			fmt.Println(bucket)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return list
}

func (this *Client) SelectBucket(name string) error {
	this.Lock()
	defer this.Unlock()
	bucket, err := this.Client.Bucket(name)
	if err != nil {
		return util.NewErrorF("select bucket err,err:%s", err)
	}
	this.selectedBucket = bucket
	return nil
}

func (this *Client) ListFile() ([]string, error) {
	if this.selectedBucket == nil {
		return nil, util.NewErrorF("not selected bucket")
	}
	bucket := this.GetSelectedBucket()
	// 列举所有文件。
	var result []string
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			return nil, util.NewErrorF("list file err,err:%s", err)
		}
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			result = append(result, object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return result, nil
}

func (this *Client) ListFiles(dir string) {
	bucket := this.GetSelectedBucket()
	marker := oss.Marker("")
	prefix := oss.Prefix(dir)
	for {
		lor, err := bucket.ListObjects(marker, prefix)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}

		for _, object := range lor.Objects {
			fmt.Printf("%s %d %s\n", object.LastModified, object.Size, object.Key)
		}

		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
}

func (this *Client) ListDirs(dir string) ([]string, error) {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	bucket := this.GetSelectedBucket()
	marker := oss.Marker("")
	prefix := oss.Prefix(dir)
	result := make([]string, 0)
	for {
		lor, err := bucket.ListObjects(marker, prefix, oss.Delimiter("/"))
		if err != nil {
			return nil, util.NewErrorF("list dir err,err:%s", err)
		}
		for _, dirName := range lor.CommonPrefixes {
			result = append(result, strings.TrimPrefix(dirName, dir))
		}
		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
	return result, nil
}

func (this *Client) GetMod(fileName string) {
	bucket := this.GetSelectedBucket()
	// 获取文件的访问权限。
	aclRes, err := bucket.GetObjectACL(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("Object ACL:", aclRes.ACL)
}

func (this *Client) SetMod(fileName string, level int) error {
	bucket := this.GetSelectedBucket()
	var aclType oss.ACLType
	switch level {
	case 1:
		aclType = oss.ACLDefault
	case 2:
		aclType = oss.ACLPublicReadWrite
	case 3:
		aclType = oss.ACLPublicRead
	}
	// 获取文件的访问权限。
	err := bucket.SetObjectACL(fileName, oss.ACLPublicReadWrite)
	if err != nil {
		return util.NewErrorF("set mod err,err:%s", err)
	}
	return nil
}

//
//func (this *Client) UploadFile(filePath string) error {
//	isExist := util.FileExist(filePath)
//	if !isExist {
//		return util.NewErrorF("file not fond,file: %s", filePath)
//	}
//	bucket := this.GetAndNewBucket(this.bucketName)
//	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
//	if err != nil {
//		return util.NewErrorF("can not open file,file: %s", filePath)
//	}
//	defer file.Close()
//	err = bucket.PutObject(this.MakeUrl(util.GetFileType(filePath)), file)
//	if err != nil {
//		return util.NewErrorF("upload file err,file: %s", filePath)
//	}
//}
//
//func getFileName(filePath string) string {
//	return fmt.Sprintf("%s%s", FilePrefix, util.GetFileType(filePath))
//}

//func UploadImage(reader *multipart.File, fileName string) string {
//	defer (*reader).Close()
//	objectKey, url := GenerateURL(&fileName, ImagePrefix)
//	e := Bucket.PutObject(objectKey, *reader)
//	if e != nil {
//		handleError(e)
//	}
//	return url + Template
//}
//
//func UploadFile(reader *multipart.File, fileName string) string {
//	defer (*reader).Close()
//
//	objectKey, url := GenerateURL(&fileName, FilePrefix)
//
//	e := Bucket.PutObject(objectKey, *reader)
//
//	if e != nil {
//		handleError(e)
//	}
//
//	return url
//}

//func GenerateURL(fileName *string, path string) (string, string) {
//	prefix := GenerateUUID()
//	split := strings.Split(*fileName, ".")
//	suffix := split[len(split)-1:][0]
//	folder := time.Now().Format("2006-10-11")
//	objectKey := path + folder + "/" + prefix + "." + suffix
//	url := URLPrefix + path + folder + "/" + prefix + "." + suffix
//	return objectKey, url
//}

func GenerateUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}
