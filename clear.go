package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anthony-dong/go-tool/logger"
	"github.com/anthony-dong/go-tool/util"
)

func main() {
	run()
	//fmt.Printf("%#v", util.String2Slice("baidu"))
}

func run() {
	log := logger.NewStdOutLogger(logger.NameOp("[GO-CLEAR]"))
	firmCode := []string{util.Slice2String([]byte{0x73, 0x6f, 0x6e, 0x67, 0x67, 0x75, 0x6f}),
		util.Slice2String([]byte{0x74, 0x74, 0x79, 0x63}),
		util.Slice2String([]byte{0x74, 0x74, 0x79, 0x6f, 0x6e, 0x67, 0x63, 0x68, 0x65}),
		util.Slice2String([]byte{0x70, 0x62, 0x73}),
	}
	log.Warnf("开始检测代码: %+v", firmCode)
	allFile := make([]string, 0)
	var err error
	if allFile, err = util.GetAllFiles("./", func(fileName string) bool {
		info, err := os.Stat(fileName)
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return false
		}
		dir := filepath.Dir(fileName)
		base := filepath.Base(fileName)
		if strings.HasPrefix(dir, ".idea") {
			return false
		}
		if strings.HasPrefix(dir, ".git") {
			return false
		}
		if strings.HasPrefix(dir, "bin") {
			return false
		}
		if base == ".DS_Store" {
			return false
		}
		return true
	}); err != nil {
		return
	}
	for _, elem := range allFile {
		log.Infof("检测文件: %s", elem)
	}
	wg := sync.WaitGroup{}
	for _, elem := range allFile {
		wg.Add(1)
		go func(elem string) {
			file, err := os.Open(elem)
			defer func() {
				wg.Done()
			}()
			if err != nil {
				panic(err)
			}
			if err := util.ReadFileLine(file, func(line string) error {
				for _, code := range firmCode {
					if strings.Contains(line, code) {
						log.Errorf("发现异常, 文件名称: %s, 检测出代码: %s", elem, code)
						panic("发现异常, 需要强制中断")
					}
				}
				return nil
			}); err != nil {
				panic(elem)
			}

		}(elem)
	}

	wg.Wait()

	log.Warnf("检测代码完成: %+v", firmCode)
}

func exit() {
	os.Exit(-1)
}
