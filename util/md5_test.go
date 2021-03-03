package util

import (
	"fmt"
	"testing"
)

func TestMd5String(t *testing.T) {
	fmt.Println(Md5([]byte("hello world")))
	fmt.Println(Md5([]byte("hello world")))
}
