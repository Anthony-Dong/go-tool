package util

import (
	"testing"
)

func Test_newKeywordMap(t *testing.T) {
	t.Logf("%#v",newKeywordMap([]string{"baidu","ali"}))
}
