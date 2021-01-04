package cli

import (
	"flag"
	"fmt"
	"os"
)

const (
	GlobeExit  = -1
	helpString = `version:1.0.0
		Usage: upload ./Main.java -config=aliyun-oss-upload-config.json
		`
)

var (
	FilePath string
	help     bool
	PutType  string
)

func InitCli() {
	flag.StringVar(&PutType, "t", "default", "-t=default")
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&FilePath, "config", "aliyun-oss-upload-config.json", "配置文件位置")

	flag.Parse()

	if help {
		fmt.Println(helpString)
		flag.PrintDefaults()
		os.Exit(GlobeExit)
	}
}
