package collections

import (
	"github.com/bytedance/gopkg/collection/skipset"
	"testing"
)

func TestSet(t *testing.T) {
	l := skipset.NewString()

	l.Add("1")
	l.Add("2")
	l.Add("3")
	l.Add("4")

	l.Range(func(value string) bool {
		t.Log(value)
		return true
	})
}
