package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"
	"simple-sync/src/utils"

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

// ValidateApiKey validates an API key and returns the associated user ID
func (s *AuthService) ValidateApiKey(apiKey string) (string, error) {
	// Get all API keys and find the one that matches
	apiKeys, err := s.storage.GetAllAPIKeys()
	if err != nil {
		return "", errors.New("failed to retrieve API keys")
	}

	for _, apiKeyModel := range apiKeys {
		if bcrypt.CompareHashAndPassword([]byte(apiKeyModel.KeyHash), []byte(apiKey)) == nil {
			// Update last used timestamp (skip for test keys to improve performance)
			if apiKeyModel.Description != "Test API Key" {
				apiKeyModel.UpdateLastUsed()
				err = s.storage.UpdateAPIKey(apiKeyModel)
				if err != nil {
					// Log error but don't fail authentication
					log.Printf("failed to update API key last used: %v", err)
				}
			}
			return apiKeyModel.UserID, nil
		}
	}

	return "", errors.New("invalid API key")
}

// GenerateApiKey generates a new API key for a user
func (s *AuthService) GenerateApiKey(userID, description string) (*models.APIKey, string, error) {
	// Generate a new API key
	plainKey, err := utils.GenerateApiKey()
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
	// Verify user exists
	_, err := s.storage.GetUserById(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Invalidate any existing setup tokens for this user
	err = s.storage.InvalidateUserSetupTokens(userID)
	if err != nil {
		return nil, errors.New("failed to invalidate existing tokens")
	}

	// Generate a new token
	token, err := utils.GenerateToken()
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
	apiKey, plainKey, err := s.GenerateApiKey(setupToken.UserID, description)
	if err != nil {
		return nil, "", err
	}

	return apiKey, plainKey, nil
}
