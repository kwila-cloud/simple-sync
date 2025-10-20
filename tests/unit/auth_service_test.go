package unit

import (
	"testing"

	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSetupToken(t *testing.T) {
	store := storage.NewTestStorage(nil)
	authService := services.NewAuthService(store)

	// Test generating setup token for existing user
	setupToken, err := authService.GenerateSetupToken(storage.TestingUserId)
	assert.NoError(t, err)
	assert.NotEmpty(t, setupToken.Token)
	assert.Equal(t, storage.TestingUserId, setupToken.User)
	assert.NotNil(t, setupToken.ExpiresAt)

	// Test generating for non-existent user
	_, err = authService.GenerateSetupToken("non-existent")
	assert.Error(t, err)
}

func TestExchangeSetupToken(t *testing.T) {
	store := storage.NewTestStorage(nil)
	authService := services.NewAuthService(store)

	// Generate setup token
	setupToken, err := authService.GenerateSetupToken(storage.TestingUserId)
	assert.NoError(t, err)

	// Test valid exchange
	apiKey, plainKey, err := authService.ExchangeSetupToken(setupToken.Token, "Test Client")
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey.UUID)
	assert.NotEmpty(t, plainKey)
	assert.Equal(t, storage.TestingUserId, apiKey.User)
	assert.Equal(t, "Test Client", apiKey.Description)

	// Test invalid token
	_, _, err = authService.ExchangeSetupToken("invalid-token", "Test")
	assert.Error(t, err)

	// Test used token (exchange again)
	_, _, err = authService.ExchangeSetupToken(setupToken.Token, "Test")
	assert.Error(t, err)
}

func TestValidateApiKey(t *testing.T) {
	store := storage.NewTestStorage(nil)
	authService := services.NewAuthService(store)

	// Generate and exchange setup token to get API key
	setupToken, err := authService.GenerateSetupToken(storage.TestingUserId)
	assert.NoError(t, err)
	_, plainKey, err := authService.ExchangeSetupToken(setupToken.Token, "Test")
	assert.NoError(t, err)

	// Test valid API key
	userID, err := authService.ValidateApiKey(plainKey)
	assert.NoError(t, err)
	assert.Equal(t, storage.TestingUserId, userID)

	// Test invalid API key
	_, err = authService.ValidateApiKey("invalid-key")
	assert.Error(t, err)
}
