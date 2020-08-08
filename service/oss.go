// @date : 2020/1/13 16:37
// @author : <a href='mailto:fanhaodong516@qq.com'>Anthony</a>

package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/anthony-dong/aliyun-oss-cli/config"
	"github.com/satori/go.uuid"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

func handleError(err error) {
	log.Println("Error:", err)
	os.Exit(-1)
}

var Client *oss.Client

var Bucket *oss.Bucket

var Template string

const (
	ImagePrefix string = "image/"
	FilePrefix  string = "file/"
)

var URLPrefix string

func UploadInit() {
	endpoint := config.OSSConfig.Endpoint
	accessKeyId := config.OSSConfig.AccessKeyId
	accessKeySecret := config.OSSConfig.AccessKeySecret
	bucketName := config.OSSConfig.BucketName
	if config.OSSConfig.FastEndpoint == "" {
		URLPrefix = "https://" + bucketName + "." + config.OSSConfig.FastEndpoint + "/"
	} else {
		URLPrefix = "https://" + bucketName + "." + endpoint + "/"
	}
	if config.OSSConfig.StyleName == "" {
		Template = ""
	} else {
		Template = "?" + strings.TrimSpace(config.OSSConfig.StyleName)
	}

	var err error
	Client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	Bucket, err = Client.Bucket(bucketName)
	if err != nil {
		Client.CreateBucket(bucketName)
		log.Println("创建Bucket 成功 , 请您修改bucket一些配置")
		handleError(err)
	}

}

func UploadImage(reader *multipart.File, fileName string) string {
	defer (*reader).Close()

	objectKey, url := GenerateURL(&fileName, ImagePrefix)

	e := Bucket.PutObject(objectKey, *reader)
	if e != nil {
		handleError(e)
	}

	return url + Template
}

func UploadFile(reader *multipart.File, fileName string) string {
	defer (*reader).Close()

	objectKey, url := GenerateURL(&fileName, FilePrefix)

	e := Bucket.PutObject(objectKey, *reader)

	if e != nil {
		handleError(e)
	}

	return url
}

func GenerateURL(fileName *string, path string) (string, string) {

	prefix := GenerateUUID()

	split := strings.Split(*fileName, ".")

	suffix := split[len(split)-1:][0]

	folder := time.Now().Format("2006-10-11")
	// path
	objectKey := path + folder + "/" + prefix + "." + suffix

	url := URLPrefix + path + folder + "/" + prefix + "." + suffix
	return objectKey, url
}

func GenerateUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}
