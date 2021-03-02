package hexo

import (
	"fmt"
	"testing"

	"github.com/gosimple/slug"
)

func TestGetFileMd5(t *testing.T) {
	fmt.Println(GetFileMd5("/Users/fanhaodong/note/1714.jpg"))
}

func Test_getOriginDir(t *testing.T) {
	list, err := checkFile("/Users/fanhaodong/note/note")
	if  err != nil {
		t.Fatal(err)
	}
	fmt.Println(list)
}

func TestSlug(t *testing.T) {
	fmt.Println(slug.MakeLang("Go 开发埋点问题.md", "cn"))
}

func TestChnageFileName(t *testing.T) {
	fmt.Println(ChangeFileName("/Users/fanhaodong/note/note/Golang/Go 开发埋点问题.md"))
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFileMd5("/Users/fanhaodong/note/1714.jpg")
	}
}
