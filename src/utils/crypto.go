package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts data using AES-256-GCM
func Encrypt(plaintext []byte, key []byte) (string, error) {
	if len(plaintext) == 0 {
		return "", errors.New("empty plaintext input")
	}
	if len(key) == 0 {
		return "", errors.New("empty key input")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-256-GCM
func Decrypt(ciphertext string, key []byte) ([]byte, error) {
	if ciphertext == "" {
		return nil, errors.New("empty ciphertext input")
	}
	if len(key) == 0 {
		return nil, errors.New("empty key input")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

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

// ExtractTokenFromHeader extracts the token from "Bearer <token>" Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[7:], nil
}
