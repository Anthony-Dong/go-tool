package gfile

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
)

/**
判断文件是否存在.
*/
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

/**
获取文件路径.
*/
func GetFileDir(path string) string {
	return filepath.Dir(path)
}

/**
获取文件的绝对路径.
*/
func GetFileAbsPath(path string) (string, error) {
	if !Exist(path) {
		return "", errors.Errorf("the file: %s not exist", path)
	}
	return filepath.Abs(path)
}

/**
获取文件名和文件后缀.
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
获取 home path.
*/
func HomePath() string {
	path, _ := exec.LookPath(os.Args[0])
	path, _ = filepath.Abs(path)
	return filepath.Dir(path)
}

func ReadFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func WriteFileBody(filename string, body []byte) error {
	if body == nil {
		return nil
	}
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(body); err != nil {
		return err
	}
	return nil
}

func ReadFileLine(file io.Reader, foo func(line string) error) error {
	reader := bufio.NewReader(file)
	for {
		lines, isEOF, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Trace(err)
		}
		if isEOF {
			break
		}
		if err := foo(string(lines)); err != nil {
			return err
		}
	}
	return nil
}
