package util

import (
	"github.com/google/uuid"
	"strings"
)

func GenerateUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
