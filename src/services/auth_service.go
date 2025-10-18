package services

import (
	"encoding/base64"
	"errors"
	"log"
	"sync"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"
	"simple-sync/src/utils"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	storage         storage.Storage
	validationMutex sync.Mutex
}

// NewAuthService creates a new auth service
func NewAuthService(storage storage.Storage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

// ValidateApiKey validates an API key and returns the associated user ID
func (s *AuthService) ValidateApiKey(apiKey string) (string, error) {
	s.validationMutex.Lock()
	defer s.validationMutex.Unlock()

	// Validate API key format before expensive operations
	if len(apiKey) < 3 || apiKey[:3] != "sk_" {
		return "", errors.New("invalid API key format")
	}

	// Check if the remaining part is valid base64 (try with padding since keys are truncated)
	base64Part := apiKey[3:]
	// Try decoding as-is first
	if _, err := base64.StdEncoding.DecodeString(base64Part); err != nil {
		// Try with padding added (generated keys are truncated to 43 chars)
		if _, err := base64.StdEncoding.DecodeString(base64Part + "="); err != nil {
			return "", errors.New("invalid API key format")
		}
	}

	// Get all API keys and find the one that matches
	apiKeys, err := s.storage.GetAllApiKeys()
	if err != nil {
		return "", errors.New("failed to retrieve API keys")
	}

	for _, apiKeyModel := range apiKeys {
		if bcrypt.CompareHashAndPassword([]byte(apiKeyModel.KeyHash), []byte(apiKey)) == nil {
			// Update last used timestamp asynchronously to avoid blocking authentication
			// Create a copy to avoid race conditions (manual copy to avoid mutex issues)
			keyCopy := &models.APIKey{
				UUID:        apiKeyModel.UUID,
				UserID:      apiKeyModel.UserID,
				KeyHash:     apiKeyModel.KeyHash,
				CreatedAt:   apiKeyModel.CreatedAt,
				LastUsedAt:  apiKeyModel.LastUsedAt,
				Description: apiKeyModel.Description,
			}
			go func() {
				keyCopy.UpdateLastUsed()
				if err := s.storage.UpdateApiKey(keyCopy); err != nil {
					log.Printf("failed to update API key last used: %v", err)
				}
			}()
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

	// Hash the API key for authentication
	keyHash, err := bcrypt.GenerateFromPassword([]byte(plainKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.New("failed to hash API key")
	}

	// Create API key model
	apiKey := models.NewAPIKey(userID, string(keyHash), description)

	// Store the API key
	err = s.storage.CreateApiKey(apiKey)
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
