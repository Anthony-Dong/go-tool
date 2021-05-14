package ghttp

import (
	"testing"
)

func TestNewHeader(t *testing.T) {
	value := "content-type  :  application/json:1111"
	headers, errr := NewHeader([]string{value})
	if errr != nil {
		t.Fatal(errr)
	}
	t.Log(headers)
}
