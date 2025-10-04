package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestAPIKeyGeneration(t *testing.T) {
	store := storage.NewTestStorage(nil)
	authService := services.NewAuthService(store)

	// Generate API key
	apiKey, plainKey, err := authService.GenerateApiKey(storage.TestingUserId, "Test Client")
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)
	assert.NotEmpty(t, plainKey)
	assert.Contains(t, plainKey, "sk_")
	assert.Equal(t, storage.TestingUserId, apiKey.UserID)
	assert.Equal(t, "Test Client", apiKey.Description)
}

func TestSetupTokenGeneration(t *testing.T) {
	store := storage.NewTestStorage(nil)
	authService := services.NewAuthService(store)

	// Create and save a test user
	user := &models.User{Id: storage.TestingUserId}
	err := store.SaveUser(user)
	assert.NoError(t, err)

	// Generate setup token
	token, err := authService.GenerateSetupToken(storage.TestingUserId)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, storage.TestingUserId, token.UserID)
	assert.Regexp(t, `^[A-Z0-9]{4}-[A-Z0-9]{4}$`, token.Token)
	assert.True(t, token.UsedAt.IsZero())
	assert.True(t, token.ExpiresAt.After(token.ExpiresAt.Add(-25*time.Hour))) // Expires in ~24 hours
}

func TestSetupTokenValidation(t *testing.T) {
	store := storage.NewTestStorage(nil)
	authService := services.NewAuthService(store)

	// Create and save a test user
	user := &models.User{Id: storage.TestingUserId}
	err := store.SaveUser(user)
	assert.NoError(t, err)

	// Generate setup token
	token, err := authService.GenerateSetupToken(storage.TestingUserId)
	assert.NoError(t, err)

	// Test valid token
	assert.True(t, token.IsValid())

	// Test expired token (simulate)
	token.ExpiresAt = token.ExpiresAt.Add(-48 * time.Hour)
	assert.False(t, token.IsValid())

	// Test used token
	token.ExpiresAt = token.ExpiresAt.Add(48 * time.Hour) // Reset expiry
	token.MarkUsed()
	assert.False(t, token.IsValid())
}
