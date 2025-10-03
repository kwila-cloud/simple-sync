package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateApiKey generates a cryptographically secure random API key
func GenerateApiKey() (string, error) {
	keyBytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(keyBytes); err != nil {
		return "", err
	}
	return "sk_" + base64.StdEncoding.EncodeToString(keyBytes)[:43], nil
}

// GenerateToken generates a random 8-character token with hyphen
func GenerateToken() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenBytes := make([]byte, 8)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := make([]byte, 9) // 4 + 1 + 4 = 9
	for i := range 4 {
		token[i] = charset[tokenBytes[i]%byte(len(charset))]
	}
	token[4] = '-'
	for i := 4; i < 8; i++ {
		token[i+1] = charset[tokenBytes[i]%byte(len(charset))]
	}

	return string(token), nil
}
