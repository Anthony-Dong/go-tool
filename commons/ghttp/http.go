package ghttp

import (
	"net/http"
	"strings"

	"github.com/juju/errors"
)

func NewHeader(headers []string) (http.Header, error) {
	result := http.Header{}
	for _, value := range headers {
		split := strings.SplitN(value, ":", 2)
		if split == nil || len(split) != 2 {
			return nil, errors.Errorf("the header %s not support", value)
		}
		result.Set(strings.TrimSpace(split[0]), strings.TrimSpace(split[1]))
	}
	return result, nil
}
