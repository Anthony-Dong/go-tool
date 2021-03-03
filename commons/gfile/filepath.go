package gfile

import (
	"fmt"
	"path/filepath"
	"strings"
)

// file-name 文件
// path   file-name的位置
func GetFileRelativePath(fileName string, path string) (string, error) {
	var err error
	if fileName, err = AbsPath(fileName); err != nil {
		return "", err
	}
	if path, err = AbsPath(path); err != nil {
		return "", err
	}
	// 没有前缀说明不在目录
	if !strings.HasPrefix(fileName, path) {
		return "", fmt.Errorf("the file %v not in path %v", fileName, path)
	}
	relativePath := strings.TrimPrefix(fileName, path)
	relativePath = filepath.Clean(relativePath)
	if strings.HasPrefix(relativePath, string(filepath.Separator)) {
		return filepath.Clean(strings.TrimPrefix(relativePath, string(filepath.Separator))), nil
	}
	return relativePath, nil
}

func AbsPath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}
