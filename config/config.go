package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

var (
	ossConfig = &Config{}
	filePath  string
)

func GetOssConfig() *Config {
	return ossConfig
}

func init() {
	flag.StringVar(&filePath, "-config", "aliyun-oss-upload-config.json", "配置文件位置")
	flag.Parse()
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
	err = json.NewDecoder(file).Decode(ossConfig)
	if err != nil {
		jc, _ := json.Marshal(ossConfig)
		logger.FatalF("读取配置文件失败,文件路径:%s,json规则:%s,config:%s", filePath, jc)
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
