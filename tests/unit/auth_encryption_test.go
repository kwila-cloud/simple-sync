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
	// Create auth service with memory storage
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-encryption-key-32-bytes-123", store)

	// Generate API key
	apiKey, plainKey, err := authService.GenerateApiKey("user-123", "Test Client")
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)
	assert.NotEmpty(t, plainKey)
	assert.Contains(t, plainKey, "sk_")
	assert.Equal(t, "user-123", apiKey.UserID)
	assert.Equal(t, "Test Client", apiKey.Description)
}

func TestSetupTokenGeneration(t *testing.T) {
	// Create auth service with memory storage
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-encryption-key-32-chars-long", store)

	// Create and save a test user
	user := &models.User{Id: "user-123"}
	err := store.SaveUser(user)
	assert.NoError(t, err)

	// Generate setup token
	token, err := authService.GenerateSetupToken("user-123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "user-123", token.UserID)
	assert.Regexp(t, `^[A-Z0-9]{4}-[A-Z0-9]{4}$`, token.Token)
	assert.False(t, token.Used)
	assert.True(t, token.ExpiresAt.After(token.ExpiresAt.Add(-25*time.Hour))) // Expires in ~24 hours
}

func TestSetupTokenValidation(t *testing.T) {
	// Create auth service with memory storage
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-encryption-key-32-chars-long", store)

	// Create and save a test user
	user := &models.User{Id: "user-123"}
	err := store.SaveUser(user)
	assert.NoError(t, err)

	// Generate setup token
	token, err := authService.GenerateSetupToken("user-123")
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
