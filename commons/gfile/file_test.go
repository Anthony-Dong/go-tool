package gfile

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"
)

func TestFileExist(t *testing.T) {
	fmt.Println(filepath.Ext("/data/open"))
	fmt.Println(path.Ext("/data/open.name"))
	fmt.Println(path.Base("/data/open.name"))
	fmt.Println(filepath.Clean("/data///open.name"))
	fmt.Println(path.Clean("/data///open.name"))
}


func TestGetFilePrefixAndSuffix(t *testing.T) {
	fmt.Println(GetFilePrefixAndSuffix("/darta/a.png"))
}
