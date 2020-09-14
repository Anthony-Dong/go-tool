package util

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func GenerateUUID() string {
	u1 := uuid.NewV4()
	return strings.ReplaceAll(u1.String(), "-", "")
}
