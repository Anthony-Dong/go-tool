package collections

import (
	"testing"
)

func TestContainsString(t *testing.T) {
	t.Log(ContainsString([]string{"elem", "elem1"}, "elem"))
}

/**
➜  go-tool git:(master) ✗ go test -run=none  -bench=^BenchmarkContainsString$ -benchmem ./commons/collections
goos: darwin
goarch: amd64
pkg: github.com/anthony-dong/go-tool/commons/collections
BenchmarkContainsString-12       7096160               160 ns/op              64 B/op          3 allocs/op
PASS
ok      github.com/anthony-dong/go-tool/commons/collections     1.316s

*/
func BenchmarkContainsString(b *testing.B) {
	elems := []string{"elem", "elem1", "elem2", "elem3", "elem4"}
	for i := 0; i < b.N; i++ {
		if !ContainsString(elems, "elem") {
			b.Fatal("error")
		}
	}
}

/**
➜  go-tool git:(master) ✗ go test -run=none  -bench=^BenchmarkContainsStringV2$ -benchmem ./commons/collections
goos: darwin
goarch: amd64
pkg: github.com/anthony-dong/go-tool/commons/collections
BenchmarkContainsStringV2-12            248700314                4.71 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/anthony-dong/go-tool/commons/collections     1.668s
*/
func BenchmarkContainsStringV2(b *testing.B) {
	elems := []string{"elem", "elem1", "elem2", "elem3", "elem4"}
	ContainsString := func(slice []string, str string) bool {
		for _, elem := range slice {
			if elem == str {
				return true
			}
		}
		return false
	}
	for i := 0; i < b.N; i++ {
		if !ContainsString(elems, "elem") {
			b.Fatal("error")
		}
	}
}
