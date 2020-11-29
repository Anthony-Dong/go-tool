package main

import (
	"fmt"
	"github.com/anthony-dong/upload-file-cli/cli"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/anthony-dong/upload-file-cli/config"
	"github.com/anthony-dong/upload-file-cli/file"
	"github.com/anthony-dong/upload-file-cli/logger"
	"github.com/anthony-dong/upload-file-cli/util"
)

func main() {
	Main()
	os.Exit(cli.GlobeExit)
}
func Main() {

	// cli
	cli.InitCli()

	// config
	config.InitConfig()

	// 获取上传文件的路径
	filePath, err := getUploadFile()

	// 文件信息设置
	year, month, day := time.Now().Date()
	fileInfo := &file.UploadFileInfo{
		LocalPath: filePath,
		Prefix:    fmt.Sprintf("%d/%d-%d", year, uint8(month), day),
		Suffix:    filepath.Ext(filePath),
		FileName:  util.GenerateUUID(),
	}
	// 打开文件
	uploadFile, err := os.Open(fileInfo.LocalPath)
	if err != nil {
		logger.FatalF("open file err,err=%v", err)
	}
	defer uploadFile.Close()

	ossConfig := config.GetOssConfigByName(cli.PutType)

	//client
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret, func(client *oss.Client) {
		client.Config.Timeout = 30
	})
	if err != nil {
		logger.FatalF("open oss err,err=%s", err)
	}

	// bucket
	bucket, err := client.Bucket(ossConfig.Bucket)
	if err != nil {
		logger.FatalF("select oss err,err=%s", err)
	}

	// file
	storageType := oss.ObjectStorageClass(oss.StorageStandard)
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)
	putPath := fileInfo.GetPutPath(ossConfig)
	err = bucket.PutObject(putPath, uploadFile, storageType, objectAcl)
	if err != nil {
		logger.FatalF("upload file err,err=%s,path=%s,err=%+v", putPath, err)
	}
	fmt.Println(fileInfo.GetOSSUrl(ossConfig))

}

func getUploadFile() (string, error) {
	filePath := handlerArgs()
	if filePath == "" {
		logger.FatalF("文件参数不存在", filePath)
	}
	if !util.FileExist(filePath) {
		logger.FatalF("文件:%s不存在", filePath)
	}
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		logger.FatalF("获取文件的绝对路径失败，path=%v,err=%v", filePath, err)
	}
	return filePath, err
}

func handlerArgs() string {
	for index, elem := range os.Args {
		if index == 0 {
			continue
		}
		if strings.HasPrefix(elem, "-") {
			continue
		}
		return elem
	}
	return ""
}
