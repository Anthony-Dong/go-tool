package digest

import (
	"fmt"
	"testing"
)

func TestMd5String(t *testing.T) {
	fmt.Println(Md5([]byte("hello world")))
	fmt.Println(Md5([]byte("hello world")))
}

/**
➜  go-tool git:(master) ✗ go test -run=none -bench=BenchmarkMd5 -benchmem ./commons/codec/digest
goos: darwin
goarch: amd64
pkg: github.com/anthony-dong/go-tool/commons/codec/digest
BenchmarkMd5-12           797011              1502 ns/op             176 B/op          4 allocs/op
PASS
*/
func BenchmarkMd5(b *testing.B) {
	str := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		Md5(str)
	}
}

func BenchmarkMd5String(b *testing.B) {
	slice := make([]byte, 1024)
	str := string(slice)
	for i := 0; i < b.N; i++ {
		Md5String(str)
	}
}
