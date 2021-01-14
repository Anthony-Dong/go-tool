package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
)

var (
	Abs = filepath.Abs
)

/**
获取全部的文件
*/
type FilterFile func(fileName string) bool

func GetAllFiles(dirPth string, filter FilterFile) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPth, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filter(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

/**
替换文件中的字符
*/
func ReplaceFileContent(old []string, new string, fileName string) error {
	if len(old) == 0 {
		return nil
	}
	readFile, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return errors.Trace(err)
	}
	defer readFile.Close()

	info, err := readFile.Stat()
	if err != nil {
		return errors.Trace(err)
	}
	fileMod := info.Mode()

	reader := bufio.NewReader(readFile)
	fileLines := make([]string, 0)
	count := 0
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

		line := string(lines)
		for _, elem := range old {
			if strings.Contains(line, elem) {
				line = strings.ReplaceAll(line, elem, new)
				count++
			}
		}
		fileLines = append(fileLines, line)
	}
	if count == 0 {
		return nil
	}
	writeFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, fileMod)
	if err != nil {
		return errors.Trace(err)
	}
	defer writeFile.Close()
	for _, elem := range fileLines {
		_, err := fmt.Fprintln(writeFile, elem)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}
