package config

import (
	"encoding/json"
	"flag"
	"github.com/anthony-dong/aliyun-oss-cli/logger"
	"github.com/anthony-dong/aliyun-oss-cli/util"
	"os"
	"path/filepath"
)

var Json = `{
  "serverPort":"服务器IP地址 例如, 9999",
  "serverHost":"你的服务器地址, 如果是远程服务器,请设置服务器IP,如果是本地 , 直接写localhost即刻",
  "serverAuth":"<请求地址必须是 : http://localhost:9999/?auth=password ,此时你需要在这里配置的就是password , 你的校验码,比如写123456>",
  "accessKeyId": "<accessKeyId>",
  "accessKeySecret": "<accessKeySecret>",
  "bucketName": "<bucketName>",
  "endpoint": "oss-accelerate.aliyuncs.com,注意golang这个上传只能使用endpoint,不能使用fastEndpoint,所以注意",
  "styleName": "x-oss-process=style/template01",
  "fastEndpoint":"oss-accelerate.aliyuncs.com,加速通道,使用fastEndpoint,没有开通就是使用普通的endpoint"
}`

type Config struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
	StyleName       string `json:"styleName"`
	FastEndpoint    string `json:"fastEndpoint"`
	ServerPort      string `json:"serverPort"`
	ServerHost      string `json:"serverHost"`
	ServerAuth      string `json:"serverAuth"`
}

var (
	ossConfig = &Config{}
	filePath  string
)

func GetOssConfig() *Config {
	return ossConfig
}

func init() {
	flag.StringVar(&filePath, "-c", "config.json", "配置文件位置")
	flag.Parse()
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		logger.FatalF("get file path err,err=%+v", err)
	}
	filePath = absPath

	isExist := util.FileExist(filePath)
	if !isExist {
		logger.FatalF("请创建文件%s", filePath)
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logger.FatalF("读取配置文件:%s失败", filePath)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(ossConfig)
	if err != nil {
		logger.FatalF("读取配置文件:%s失败", filePath)
	}
}
