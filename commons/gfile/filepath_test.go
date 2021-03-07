package gfile

import (
	"fmt"
	"testing"
)

func TestGetFileRelativePath(t *testing.T) {
	fmt.Println(GetFileRelativePath("/data/log/test/a.log", "/data"))
}

func BenchmarkGetFileRelativePath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFileRelativePath("/data/log", "/data")
	}
}
