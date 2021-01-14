package util

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
