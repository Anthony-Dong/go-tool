package collections

import (
	"fmt"
	"testing"
)

func TestGetKeys(t *testing.T) {
	strings, err := GetMapKeysToString(map[string]struct{}{"k1": {}, "k2": {}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(strings)
}

func TestGetKeys2(t *testing.T) {
	strings, err := GetMapKeysToString(map[interface{}]struct{}{1: {}, "k2": {}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(strings)
}

func TestToCliMultiDescString(t *testing.T) {
	t.Log(ToCliMultiDescString([]string{"k1"}))
	t.Log(ToCliMultiDescString([]string{"k1", "k2"}))
}

func TestSplitStringSlice(t *testing.T) {
	fmt.Println(SplitStringSlice([]string{"1", "2", "3"}, 2))
	fmt.Println(SplitStringSlice([]string{"1", "2"}, 2))
	fmt.Println(SplitStringSlice([]string{"1"}, 2))
}
