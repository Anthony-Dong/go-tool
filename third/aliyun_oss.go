package third

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/anthony-dong/go-tool/util"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/juju/errors"
)

type OssConfigs map[string]OssConfig

func (o OssConfigs) GetConfig(configName string) *OssConfig {
	config, isExist := o[configName]
	if isExist {
		return &config
	}
	return nil
}

type OssConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	UrlEndpoint     string `json:"url_endpoint"`
	Bucket          string `json:"bucket"`
	PathPrefix      string `json:"path_prefix"`
}

type OssUploadFile struct {
	LocalFile  string `json:"local_file"` // 本地文件
	SaveDir    string `json:"save_dir"`   // 保存到远程的地址
	FilePrefix string `json:"file_name"`  // 文件名称
	FileSuffix string `json:"file_type"`  // 文件类型名称
}

// image/2019-08-29/38564c69-85ba-4415-93d8-cb05e783c4b6.jpg
func (f *OssUploadFile) GetPutPath(config *OssConfig) string {
	return filepath.Join(config.PathPrefix, f.SaveDir, fmt.Sprintf("%s%s", f.FilePrefix, f.FileSuffix))
}

// https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/d21baa6d76a14aa8b70db1c033891990.png
func (f *OssUploadFile) GetOSSUrl(config *OssConfig) string {
	path := f.GetPutPath(config)
	return fmt.Sprintf("https://%s/%s", config.UrlEndpoint, path)
}

/**
new Bucket
*/
func NewBucket(ossConfig *OssConfig) (*oss.Bucket, error) {
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret, func(client *oss.Client) {
		client.Config.Timeout = 5
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	// bucket
	bucket, err := client.Bucket(ossConfig.Bucket)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return bucket, nil
}

/**
上传文件
*/
func (f *OssUploadFile) PutFile(bucket *oss.Bucket, ossConfig *OssConfig) error {
	file, err := os.Open(f.LocalFile)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() {
		_ = file.Close()
	}()
	return errors.Trace(bucket.PutObject(f.GetPutPath(ossConfig), file, oss.ObjectStorageClass(oss.StorageStandard), oss.ObjectACL(oss.ACLPublicRead)))
}

/**
获取配置文件
*/
func GetOssConfig(configFile string) (OssConfigs, error) {
	filePath, err := util.GetFilePath(configFile)
	if err != nil {
		return nil, errors.Trace(err)
	}
	_file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer _file.Close()
	body, err := ioutil.ReadAll(_file)
	if err != nil {
		return nil, errors.Trace(err)
	}
	var (
		multiOssConfig = map[string]OssConfig{}
		ossConfig      = &OssConfig{}
		isSingle       = false
	)
	// 原来是单个配置文件
	{
		err = json.Unmarshal(body, ossConfig)
		if err != nil {
			return nil, err
		}
		if ossConfig.AccessKeyId != "" {
			isSingle = true
		}
	}
	//现在需要支持多个
	{
		err = json.Unmarshal(body, &multiOssConfig)
		if !isSingle && err != nil {
			return nil, err
		}
	}
	if isSingle {
		multiOssConfig["default"] = *ossConfig
	}
	return multiOssConfig, nil
}
