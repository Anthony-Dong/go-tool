package digest

import (
	"crypto/md5"
	"fmt"

	"github.com/anthony-dong/go-tool/commons/gstring"
)

func Md5String(data string) string {
	return Md5(gstring.String2Slice(data))
}

func Md5(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
