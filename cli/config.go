package cli

import (
	"flag"
	"fmt"
	"os"
)

const (
	GlobeExit = -1
)

var (
	FilePath string
	help     bool
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&FilePath, "config", "aliyun-oss-upload-config.json", "配置文件位置")
	flag.Parse()

	if help {
		fmt.Println(`version:1.0.0
Usage: upload ./Main.java -config=aliyun-oss-upload-config.json
`)
		flag.PrintDefaults()
		os.Exit(GlobeExit)
	}
}
