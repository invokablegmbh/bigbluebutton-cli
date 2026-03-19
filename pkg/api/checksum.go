package api

import (
	"crypto/sha256"
	"encoding/hex"
)

// Checksum calculates the BigBlueButton API checksum.
// It concatenates apiCall + queryString + secret and returns the SHA-256 hex digest.
func Checksum(apiCall, queryString, secret string) string {
	h := sha256.New()
	h.Write([]byte(apiCall + queryString + secret))
	return hex.EncodeToString(h.Sum(nil))
}
