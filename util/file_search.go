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

func WriteFile(writer io.Writer, body []string) error {
	if writer == nil || body == nil {
		return errors.New("params is nil")
	}
	var line = []byte{'\n'}
	for _, elem := range body {
		if _, err := writer.Write(String2Slice(elem)); err != nil {
			return err
		}
		if _, err := writer.Write(line); err != nil {
			return err
		}
	}
	return nil
}

// V2
func ReplaceFileContentV2(keyword []string, file io.Reader) ([]string, bool, error) {
	if file == nil {
		return nil, false, errors.New("file is nil")
	}
	if keyword == nil || len(keyword) == 0 {
		return nil, false, nil
	}
	count := 0
	keywordMap := newKeywordMap(keyword)
	reader := bufio.NewReader(file)
	fileLines := make([]string, 0)
	hasContent := false
	for {
		lines, isEOF, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, false, errors.Trace(err)
		}
		if isEOF {
			break
		}
		line := string(lines)
		if line == "" && !hasContent {
			continue
		}
		hasContent = true
		for _, elem := range keyword {
			if strings.Contains(line, elem) {
				line = strings.ReplaceAll(line, elem, keywordMap[elem])
				count++
			}
		}
		fileLines = append(fileLines, line)
	}
	return fileLines, count > 0, nil
}

func newKeywordMap(keyword []string) map[string]string {
	if keyword == nil || len(keyword) == 0 {
		return nil
	}
	builder := strings.Builder{}
	result := make(map[string]string, len(keyword))
	for _, elem := range keyword {
		strLen := len(elem)
		if strLen == 0 {
			continue
		}
		builder.Reset()
		for x := strLen - 1; x >= 0; x-- {
			builder.WriteByte(elem[x])
		}
		result[elem] = builder.String()
	}
	return result
}
