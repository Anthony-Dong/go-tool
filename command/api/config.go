package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/anthony-dong/go-tool/commons/gfile"
	logger2 "github.com/anthony-dong/go-tool/commons/logger"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"

	"github.com/urfave/cli/v2"
)

const (
	logLevelFlag = "log-level"
	configFlag   = "config"
)

var (
	GlobalFlag = []cli.Flag{
		&cli.StringFlag{
			Name:     logLevelFlag,
			Usage:    fmt.Sprintf("Set the logging level (%s)", logger2.LogLevelToString()),
			Required: false,
			Value:    "debug",
		},
		&cli.StringFlag{
			Name:     configFlag,
			Usage:    "Location of client config files",
			Required: false,
			Value:    filepath.Join(gfile.HomePath(), ".go-tool.json"),
		},
	}
)

type CommonConfig struct {
	Config   string `json:"config"`
	LogLevel string `json:"log-level"`
}

func GetCommonConfig(context *cli.Context) CommonConfig {
	return CommonConfig{
		LogLevel: context.String(logLevelFlag),
		Config:   context.String(configFlag),
	}
}

func (c CommonConfig) ReadConfig(tag string) ([]byte, error) {
	file, err := os.Open(c.Config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer file.Close()
	all, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Trace(err)
	}
	result := gjson.GetBytes(all, tag)
	if result.Exists() {
		return []byte(result.String()), nil
	}
	return nil, errors.Errorf("config not found %s tag", tag)
}
