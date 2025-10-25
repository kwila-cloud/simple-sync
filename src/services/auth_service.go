package services

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	apperrors "simple-sync/src/errors"
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

	// Validate API key format minimally before expensive operations
	if len(apiKey) < 4 || apiKey[:3] != "sk_" {
		return "", apperrors.ErrInvalidApiKeyFormat
	}

	// Get all API keys and find the one that matches
	apiKeys, err := s.storage.GetAllApiKeys()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve API keys: %w", err)
	}

	for _, apiKeyModel := range apiKeys {
		if bcrypt.CompareHashAndPassword([]byte(apiKeyModel.KeyHash), []byte(apiKey)) == nil {
			// Update last used timestamp asynchronously to avoid blocking authentication
			// Create a copy to avoid race conditions (manual copy to avoid mutex issues)
			keyCopy := &models.ApiKey{
				UUID:        apiKeyModel.UUID,
				User:        apiKeyModel.User,
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
			return apiKeyModel.User, nil
		}
	}

	return "", apperrors.ErrInvalidApiKey
}

// GenerateApiKey generates a new API key for a user
func (s *AuthService) GenerateApiKey(userID, description string) (*models.ApiKey, string, error) {
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
	apiKey := models.NewApiKey(userID, string(keyHash), description)

	// Store the API key
	err = s.storage.AddApiKey(apiKey)
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
		if err == storage.ErrNotFound {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
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
	err = s.storage.AddSetupToken(setupToken)
	if err != nil {
		return nil, errors.New("failed to store setup token")
	}

	return setupToken, nil
}

// ExchangeSetupToken exchanges a setup token for an API key
func (s *AuthService) ExchangeSetupToken(token, description string) (*models.ApiKey, string, error) {
	// Get the setup token
	setupToken, err := s.storage.GetSetupToken(token)
	if err != nil {
		return nil, "", apperrors.ErrInvalidSetupToken
	}

	// Validate the token
	if !setupToken.IsValid() {
		return nil, "", apperrors.ErrSetupTokenExpired
	}

	// Mark the token as used
	setupToken.MarkUsed()
	err = s.storage.UpdateSetupToken(setupToken)
	if err != nil {
		return nil, "", errors.New("failed to update setup token")
	}

	// Generate API key for the user
	apiKey, plainKey, err := s.GenerateApiKey(setupToken.User, description)
	if err != nil {
		return nil, "", err
	}

	return apiKey, plainKey, nil
}
