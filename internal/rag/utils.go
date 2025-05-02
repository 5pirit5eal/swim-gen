package rag

import (
	"crypto/sha256"
	"encoding/hex"
)

// Calculates a hash for the given slice of strings
// GenerateHash generates a SHA256 hash of len 64 from a given string.
func GenerateHash(inputs ...string) string {
	hasher := sha256.New()
	bytes := []byte{}
	for _, input := range inputs {
		bytes = append(bytes, []byte(input)...)
	}
	hasher.Write(bytes)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)[:64]
}
