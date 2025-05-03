package rag

import (
	"strings"

	"github.com/google/uuid"
)

// Calculates a uuid for the given slice of strings
func GenerateUUID(inputs ...string) string {
	combined := uuid.NewSHA1(uuid.NameSpaceOID, []byte(strings.Join(inputs, "")))

	return combined.String()
}
