package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestValidateAPIKey(t *testing.T) {
	// Setup storage and auth service
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-encryption-key-32-chars-long", store)

	// Generate an API key
	_, plainKey, err := authService.GenerateAPIKey("testuser", "Test Client")
	assert.NoError(t, err)

	// Validate the API key
	userID, err := authService.ValidateAPIKey(plainKey)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", userID)

	// Test invalid API key
	_, err = authService.ValidateAPIKey("sk_invalid123456789012345678901234567890")
	assert.Error(t, err)

	// Test non-existent API key
	_, err = authService.ValidateAPIKey("sk_nonexistent123456789012345678901234567890")
	assert.Error(t, err)
}

func TestExchangeSetupToken(t *testing.T) {
	// Setup storage and auth service
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-encryption-key-32-chars-long", store)

	// Generate a setup token
	setupToken, err := authService.GenerateSetupToken("testuser")
	assert.NoError(t, err)

	// Exchange the setup token
	apiKey, plainKey, err := authService.ExchangeSetupToken(setupToken.Token, "Mobile App")
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)
	assert.NotEmpty(t, plainKey)
	assert.Equal(t, "testuser", apiKey.UserID)
	assert.Equal(t, "Mobile App", apiKey.Description)

	// Verify the token was marked as used
	updatedToken, err := store.GetSetupToken(setupToken.Token)
	assert.NoError(t, err)
	assert.True(t, updatedToken.Used)

	// Test exchanging used token
	_, _, err = authService.ExchangeSetupToken(setupToken.Token, "Another App")
	assert.Error(t, err)

	// Test exchanging invalid token
	_, _, err = authService.ExchangeSetupToken("INVALID-TOKEN", "Test")
	assert.Error(t, err)
}

func TestAPIKeyModelValidation(t *testing.T) {
	// Test valid API key
	validKey := &models.APIKey{
		UUID:         "550e8400-e29b-41d4-a716-446655440000",
		UserID:       "testuser",
		EncryptedKey: "encrypted-data",
		KeyHash:      "hash-data",
		CreatedAt:    time.Now(),
		Description:  "Test Key",
	}
	err := validKey.Validate()
	assert.NoError(t, err)

	// Test invalid API key - missing UUID
	invalidKey := &models.APIKey{
		UserID:       "testuser",
		EncryptedKey: "encrypted-data",
		KeyHash:      "hash-data",
		CreatedAt:    time.Now(),
	}
	err = invalidKey.Validate()
	assert.Error(t, err)
}

func TestSetupTokenModelValidation(t *testing.T) {
	// Test valid setup token
	validToken := &models.SetupToken{
		Token:     "ABCD-1234",
		UserID:    "testuser",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}
	err := validToken.Validate()
	assert.NoError(t, err)
	assert.True(t, validToken.IsValid())

	// Test invalid token format
	invalidToken := &models.SetupToken{
		Token:     "INVALID-FORMAT",
		UserID:    "testuser",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}
	err = invalidToken.Validate()
	assert.Error(t, err)
}
