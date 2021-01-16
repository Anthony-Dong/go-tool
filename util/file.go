package util

import (
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
获取 home path
*/
func HomePath() string {
	path, _ := exec.LookPath(os.Args[0])
	path, _ = filepath.Abs(path)
	return filepath.Dir(path)
}
