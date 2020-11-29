package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/anthony-dong/upload-file-cli/cli"

	"github.com/anthony-dong/upload-file-cli/logger"
	"github.com/anthony-dong/upload-file-cli/util"
)

type Config struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	UrlEndpoint     string `json:"url_endpoint"`
	Bucket          string `json:"bucket"`
	PathPrefix      string `json:"path_prefix"`
}

const (
	DefaultTypeName = "default"
)

var (
	demoMultiOssConfig = map[string]Config{
		"default": {},
	}
	ossConfig      = &Config{}
	multiOssConfig = map[string]Config{}
	filePath       string

	isSingle bool
)

/**
{
  "access_key_id": "xxxx",
  "access_key_secret": "xxxx",
  "endpoint": "oss-accelerate.xxx.com",
  "url_endpoint": "xxx.oss-xxx.xxx.com",
  "bucket": "xxx",
  "path_prefix": "xxx"
}
*/

func GetOssConfig() *Config {
	return ossConfig
}

func GetOssConfigByName(name string) *Config {
	config, isExist := multiOssConfig[name]
	if !isExist {
		return nil
	}
	return &config
}

func InitConfig() {
	filePath = cli.FilePath
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		logger.FatalF("get file path err,err=%+v", err)
	}
	if util.FileExist(filePath) {
		filePath = absPath
	} else {
		path := GetRootPath()
		rootConfigPath := fmt.Sprintf("%s/%s", path, filePath)
		rootConfigPath = filepath.Clean(rootConfigPath)
		filePath = rootConfigPath
	}
	file, err := os.Open(filePath)
	if err != nil {
		logger.FatalF("读取配置文件:%s失败", filePath)
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		logger.FatalF("读取配置文件:%s失败,err:%v", filePath, err)
	}

	// 原来是单个配置文件
	{
		err = json.Unmarshal(body, ossConfig)
		if err != nil {
			jc, _ := json.Marshal(ossConfig)
			logger.FatalF("读取配置文件失败,文件路径:%s,json规则:%s,config:%s", filePath, jc)
		}

		if ossConfig.AccessKeyId != "" {
			isSingle = true
		}
	}

	//现在需要支持多个

	{
		err = json.Unmarshal(body, &multiOssConfig)
		if !isSingle && err != nil {
			jc, _ := json.Marshal(demoMultiOssConfig)
			logger.FatalF("读取配置文件失败,文件路径:%s,json规则:%s,config:%s", filePath, jc)
		}
	}

	if isSingle {
		multiOssConfig[DefaultTypeName] = *ossConfig
	}
}

func GetRootPath() string {
	curFilename := os.Args[0]
	binaryPath, err := exec.LookPath(curFilename)
	if err != nil {
		panic(err)
	}
	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		panic(err)
	}
	return filepath.Dir(binaryPath)
}
