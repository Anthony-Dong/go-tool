package digest

import "net/url"

// 编码url
func Base64Encode(str string) string {
	return url.QueryEscape(str)
}

func Base64Decode(str string) (string, error) {
	return url.QueryUnescape(str)
}
