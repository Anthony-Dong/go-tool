package third

import (
	"encoding/json"
	"fmt"
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

// image/2019-08-29/38564c69-85ba-4415-93d8-xxxxx.jpg
func (f *OssUploadFile) GetPutPath(config *OssConfig) string {
	return filepath.Join(config.PathPrefix, f.SaveDir, fmt.Sprintf("%s%s", f.FilePrefix, f.FileSuffix))
}

// https://xxxx.oss-accelerate.xxxx.com/image/2020/9-14/xxxxxx.png
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
func GetOssConfig(body []byte) (OssConfigs, error) {
	var (
		multiOssConfig = map[string]OssConfig{}
	)
	if err := json.Unmarshal(body, &multiOssConfig); err != nil {
		return nil, errors.Annotatef(err, "you should set config in you config file:%s\n", util.ToJsonString(OssConfigs{"default": OssConfig{}}))
	}
	return multiOssConfig, nil
}
