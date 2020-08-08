package util

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unsafe"
)

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func FileIsDir(fileName string) (bool, error) {
	info, err := os.Stat(fileName)
	if err != nil {
		return false, NewErrorF("can not find %s", fileName)
	}
	return info.IsDir(), nil
}
func GetFileType(fileName string) string {
	return path.Ext(fileName)
}


func Slice2String(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}

func SplitString(text []byte) []string {
	str := Slice2String(text)
	split := strings.Split(str, "\\s{1,})")
	return split
}

func NewErrorF(format string, arg ...interface{}) error {
	return fmt.Errorf(format, arg ...)
}
