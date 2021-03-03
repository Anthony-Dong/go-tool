package util

import (
	"crypto/md5"
	"fmt"
)

func Md5String(data string) string {
	return Md5(String2Slice(data))
}

func Md5(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
