package util

import (
	"unsafe"
)

/**
可以安全的操作
*/
func Slice2String(body []byte) string {
	if body == nil || len(body) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&body))
}

type _string struct {
	data string
	cap  int
}

/**
1、可以安全的转换为切片
2、转换后的切片可以进行修改操作
*/
func String2Slice(body string) []byte {
	if body == "" {
		return []byte{}
	}
	str := _string{data: body, cap: len(body)}
	return *(*[]byte)(unsafe.Pointer(&str))
}
