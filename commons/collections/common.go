package collections

import (
	"github.com/thoas/go-funk"
)

func ContainsString(str []string, elem string) bool {
	return funk.Contains(str, elem)
}
