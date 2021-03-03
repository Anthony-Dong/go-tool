package digest

import (
	"testing"
)

func TestBase64(t *testing.T) {
	t.Log(Base64Encode("测试代码"))
	t.Log(Base64Decode("%E6%B5%8B%E8%AF%95%E4%BB%A3%E7%A0%81"))
}

