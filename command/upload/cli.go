package upload

import (
	"fmt"
	"time"

	"github.com/anthony-dong/go-tool/command"
	"github.com/anthony-dong/go-tool/command/api"
	"github.com/anthony-dong/go-tool/command/log"
	"github.com/anthony-dong/go-tool/commons/codec/digest"
	"github.com/anthony-dong/go-tool/commons/codec/gjson"
	"github.com/anthony-dong/go-tool/commons/collections"
	"github.com/anthony-dong/go-tool/commons/gfile"
	"github.com/anthony-dong/go-tool/commons/gos"
	"github.com/anthony-dong/go-tool/commons/gtime"
	"github.com/anthony-dong/go-tool/commons/uuid"
	"github.com/anthony-dong/go-tool/third"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

var (
	decodeType = map[string]struct{}{
		"uuid":   {},
		"base64": {},
	}
	decodeTypeName = func() string {
		list, _ := collections.GetMapKeysToString(decodeType)
		return collections.ToCliMultiDescString(list)
	}
)

type uploadCommand struct {
	api.CommonConfig
	OssConfigType  string `json:"type"`
	File           string `json:"file"`
	FileNameDecode string `json:"decode"`
}

func NewUploadCommand() command.Command {
	return new(uploadCommand)
}

func (c *uploadCommand) InitConfig(context *cli.Context, config api.CommonConfig) (_ []byte, err error) {
	c.CommonConfig = config
	c.File, err = gfile.Abs(c.File)
	if err != nil {
		return nil, errors.Annotate(err, "获取文件绝对路径失败")
	}
	_, isExist := decodeType[c.FileNameDecode]
	if !isExist {
		return nil, errors.Errorf("decode method not found: %s", c.FileNameDecode)
	}
	return gjson.ToJsonString(c), nil
}

func (c *uploadCommand) Flag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "type",
			Aliases:     []string{"t"},
			Usage:       fmt.Sprintf("Set the upload config type ()"),
			Destination: &c.OssConfigType,
			Required:    false,
			Value:       "default",
		},
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Usage:       "Set the local path of upload file",
			Destination: &c.File,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "decode",
			Aliases:     []string{"d"},
			Usage:       fmt.Sprintf("Set the upload file name decode method (%s)", decodeTypeName()),
			Destination: &c.FileNameDecode,
			Required:    false,
			Value:       "uuid",
		},
	}
}

func (c *uploadCommand) Run(context *cli.Context) error {
	jsonConfig, err := c.ReadConfig("upload")
	if err != nil {
		return errors.Trace(err)
	}
	configs, err := third.GetOssConfig(jsonConfig)
	if err != nil {
		return errors.Trace(err)
	}
	config := configs.GetConfig(c.OssConfigType)
	if config == nil {
		return errors.New("the config is nil")
	}
	prefix, suffix := gfile.GetFilePrefixAndSuffix(c.File)
	fileInfo := third.OssUploadFile{
		LocalFile:  c.File,
		SaveDir:    time.Now().Format(gtime.FromatTime_V2),
		FilePrefix: c.getFileName(prefix),
		FileSuffix: suffix,
	}
	bucket, err := third.NewBucket(config)
	if err != nil {
		return errors.Trace(err)
	}
	if err := fileInfo.PutFile(bucket, config); err != nil {
		return errors.Trace(err)
	}
	fileUrl := fileInfo.GetOSSUrl(config)
	if log.IsInfoEnabled() {
		log.Infof("[upload] end success, url: %s", fileUrl)
	} else {
		fmt.Println(fileUrl)
		gos.ExitError()
	}
	return nil
}

func (c *uploadCommand) getFileName(fileName string) string {
	switch c.FileNameDecode {
	case "uuid":
		return uuid.GenerateUUID()
	case "base64":
		return digest.Base64Encode(fileName)
	}
	return ""
}
