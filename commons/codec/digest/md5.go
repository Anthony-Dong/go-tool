package digest

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/anthony-dong/go-tool/commons/gstring"
)

func Md5String(data string) string {
	return Md5(gstring.String2Slice(data))
}

func Md5(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}
