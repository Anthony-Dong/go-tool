package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
)

/**
判断文件是否存在
*/
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

/**
获取文件路径
*/
func GetFileDir(path string) string {
	return filepath.Dir(path)
}

/**
获取文件的绝对路径
*/
func GetFileAbsPath(path string) (string, error) {
	if !Exist(path) {
		return "", errors.Errorf("the file: %s not exist", path)
	}
	return filepath.Abs(path)
}

/**
获取文件名和文件后缀
*/
func GetFilePrefixAndSuffix(filename string) (prefix, suffix string) {
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename, ""
	}
	filename = strings.TrimSuffix(filename, ext)
	return filename, ext
}

/**
获取文件路径
*/
func GetFilePath(fileName string) (string, error) {
	// 当前执行脚本路径
	var curPath = func() string {
		path, _ := exec.LookPath(os.Args[0])
		path, _ = filepath.Abs(path)
		return filepath.Dir(path)
	}
	// 判断文件是否存在
	var exist = func(filename string) bool {
		_, err := os.Stat(filename)
		return err == nil || os.IsExist(err)
	}
	str, _ := filepath.Abs(fileName)
	if exist(str) {
		return str, nil
	}
	path1 := filepath.Clean(filepath.Join(curPath(), fileName))
	if !exist(path1) {
		return "", fmt.Errorf("file: %s not fond", path1)
	}
	return path1, nil
}
