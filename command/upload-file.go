package command

import (
	"fmt"
	"time"

	"github.com/anthony-dong/go-tool/third"
	"github.com/anthony-dong/go-tool/util"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

var (
	decodeType = map[string]struct{}{
		"uuid":   {},
		"base64": {},
	}
)

type uploadCommand struct {
	OssConfigFile  string `json:"oss_config_file"`
	OssConfigType  string `json:"oss_config_type"`
	File           string `json:"file"`
	FileNameDecode string `json:"file_name_decode"`
}

func NewUploadCommand() Command {
	return new(uploadCommand)
}

func (c *uploadCommand) InitConfig(context *cli.Context) ([]byte, error) {
	configFilePath, err := util.GetFilePath(c.OssConfigFile)
	if err != nil {
		return nil, errors.Trace(err)
	}
	c.OssConfigFile = configFilePath
	file, err := util.GetFilePath(c.File)
	if err != nil {
		return nil, errors.Trace(err)
	}
	c.File = file
	_, isExist := decodeType[c.FileNameDecode]
	if !isExist {
		return nil, errors.Errorf("decode method not found: %s", c.FileNameDecode)
	}
	return util.ToJsonString(c), nil
}

func (c *uploadCommand) Flag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "oss_config_file",
			Aliases:     []string{"c"},
			Usage:       "the aliyun oss config file",
			Destination: &c.OssConfigFile,
			Required:    false,
			Value:       "aliyun-oss-upload-config.json",
		},
		&cli.StringFlag{
			Name:        "oss_config_type",
			Aliases:     []string{"t"},
			Usage:       "the aliyun oss config type, default is default",
			Destination: &c.OssConfigType,
			Required:    false,
			Value:       "default",
		},
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Usage:       "the upload file local path",
			Destination: &c.File,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "file_name_decode",
			Aliases:     []string{"d"},
			Usage:       "the upload file name decode",
			Destination: &c.FileNameDecode,
			Required:    false,
			Value:       "uuid",
		},
	}
}

func (c *uploadCommand) Run(context *cli.Context) error {
	configs, err := third.GetOssConfig(c.OssConfigFile)
	if err != nil {
		return errors.Trace(err)
	}
	config := configs.GetConfig(c.OssConfigType)
	if config == nil {
		return util.NilError("config")
	}
	prefix, suffix := util.GetFilePrefixAndSuffix(c.File)
	file := third.OssUploadFile{
		LocalFile:  c.File,
		SaveDir:    time.Now().Format(util.FromatTime_V2),
		FilePrefix: c.getFileName(prefix),
		FileSuffix: suffix,
	}
	bucket, err := third.NewBucket(config)
	if err != nil {
		return errors.Trace(err)
	}
	if err := file.PutFile(bucket, config); err != nil {
		return errors.Trace(err)
	}
	fileUrl := file.GetOSSUrl(config)
	log.Infof("[upload] end success, url: %s", fileUrl)
	fmt.Println(fileUrl)
	return nil
}

func (c *uploadCommand) getFileName(fileName string) string {
	switch c.FileNameDecode {
	case "uuid":
		return util.GenerateUUID()
	case "base64":
		return util.Base64Encode(fileName)
	}
	return ""
}
