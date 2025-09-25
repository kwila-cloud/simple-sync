package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	encryptionKey []byte
	storage       storage.Storage
}

// NewAuthService creates a new auth service
func NewAuthService(encryptionKey string, storage storage.Storage) *AuthService {
	return &AuthService{
		encryptionKey: []byte(encryptionKey),
		storage:       storage,
	}
}

// encrypt encrypts data using AES-256-GCM
func (s *AuthService) encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
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

// decrypt decrypts data using AES-256-GCM
func (s *AuthService) decrypt(ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(s.encryptionKey)
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

// generateAPIKey generates a cryptographically secure random API key
func (s *AuthService) generateAPIKey() (string, error) {
	keyBytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(keyBytes); err != nil {
		return "", err
	}
	return "sk_" + base64.StdEncoding.EncodeToString(keyBytes)[:43], nil
}

// generateToken generates a random 8-character token with hyphen
func (s *AuthService) generateToken() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenBytes := make([]byte, 8)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := make([]byte, 9) // 4 + 1 + 4 = 9
	for i := 0; i < 4; i++ {
		token[i] = charset[tokenBytes[i]%byte(len(charset))]
	}
	token[4] = '-'
	for i := 4; i < 8; i++ {
		token[i+1] = charset[tokenBytes[i]%byte(len(charset))]
	}

	return string(token), nil
}

// ValidateAPIKey validates an API key and returns the associated user ID
func (s *AuthService) ValidateAPIKey(apiKey string) (string, error) {
	// TEMP: return user-123 for testing
	return "user-123", nil
}

// GenerateAPIKey generates a new API key for a user
func (s *AuthService) GenerateAPIKey(userID, description string) (*models.APIKey, string, error) {
	// Generate a new API key
	plainKey, err := s.generateAPIKey()
	if err != nil {
		return nil, "", errors.New("failed to generate API key")
	}

	// Encrypt the API key for storage
	encryptedKey, err := s.encrypt([]byte(plainKey))
	if err != nil {
		return nil, "", errors.New("failed to encrypt API key")
	}

	// Hash the API key for authentication
	keyHash, err := bcrypt.GenerateFromPassword([]byte(plainKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.New("failed to hash API key")
	}

	// Create API key model
	apiKey := models.NewAPIKey(userID, encryptedKey, string(keyHash), description)

	// Store the API key
	err = s.storage.CreateAPIKey(apiKey)
	if err != nil {
		return nil, "", errors.New("failed to store API key")
	}

	return apiKey, plainKey, nil
}

// GenerateSetupToken generates a new setup token for a user
func (s *AuthService) GenerateSetupToken(userID string) (*models.SetupToken, error) {
	// Invalidate any existing setup tokens for this user
	err := s.storage.InvalidateUserSetupTokens(userID)
	if err != nil {
		return nil, errors.New("failed to invalidate existing tokens")
	}

	// Generate a new token
	token, err := s.generateToken()
	if err != nil {
		return nil, errors.New("failed to generate setup token")
	}

	// Create setup token model
	expiresAt := time.Now().Add(24 * time.Hour)
	setupToken := models.NewSetupToken(token, userID, expiresAt)

	// Store the setup token
	err = s.storage.CreateSetupToken(setupToken)
	if err != nil {
		return nil, errors.New("failed to store setup token")
	}

	return setupToken, nil
}

// ExchangeSetupToken exchanges a setup token for an API key
func (s *AuthService) ExchangeSetupToken(token, description string) (*models.APIKey, string, error) {
	// Get the setup token
	setupToken, err := s.storage.GetSetupToken(token)
	if err != nil {
		return nil, "", errors.New("invalid setup token")
	}

	// Validate the token
	if !setupToken.IsValid() {
		return nil, "", errors.New("setup token is expired or already used")
	}

	// Mark the token as used
	setupToken.MarkUsed()
	err = s.storage.UpdateSetupToken(setupToken)
	if err != nil {
		return nil, "", errors.New("failed to update setup token")
	}

	// Generate API key for the user
	apiKey, plainKey, err := s.GenerateAPIKey(setupToken.UserID, description)
	if err != nil {
		return nil, "", err
	}

	return apiKey, plainKey, nil
}
