package gstring

import (
	"path/filepath"
	"testing"
)

func TestSlug(t *testing.T) {
	t.Log(Slug("Docker的网络模型 - macvlan & ipvlan"))
	t.Log(filepath.Clean("Docker的网络模型 - macvlan & ipvlan.md"))
	t.Log(filepath.Split("/home/polaris/studygolang"))
}
