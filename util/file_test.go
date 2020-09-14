package util

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

func TestName(b *testing.T) {
	fmt.Println(getParent("/"))
	fmt.Println(getParent("/a"))
	fmt.Println(getParent("a/b/c/"))
}
