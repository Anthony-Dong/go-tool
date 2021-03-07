package git

import "testing"

var (
	ugnore = CompileIgnoreLines("/bin1", "/bin2", "/bin3", "/bin4", "/bin5", "/bin6", "/bin")
)

func BenchmarkMatchesPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ugnore.MatchesPath("bin/app")
	}
}
