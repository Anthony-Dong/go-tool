package gfile

import (
	"fmt"
	"strings"
	"testing"
)

func Test_newKeywordMap(t *testing.T) {
	t.Logf("%#v", newKeywordMap([]string{"baidu", "ali"}))
}

func TestGetAllFiles(t *testing.T) {
	files, err := GetAllFiles("/Users/fanhaodong/note/note", func(fileName string) bool {
		relativePath, err := GetFileRelativePath(fileName, "/Users/fanhaodong/note/note")
		if err != nil {
			panic(err)
		}
		return strings.HasSuffix(relativePath, ".md") || strings.HasSuffix(relativePath, ".markdown")
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(len(files))
}
