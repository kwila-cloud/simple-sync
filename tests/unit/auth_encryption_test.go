package unit

import (
	"testing"
	"time"

	"simple-sync/src/services"

	"github.com/stretchr/testify/assert"
)

func TestAPIKeyGeneration(t *testing.T) {
	// Create auth service
	authService := services.NewAuthService("test-encryption-key-32-chars-long", nil)

	// Generate API key
	apiKey, plainKey, err := authService.GenerateAPIKey("testuser", "Test Client")
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)
	assert.NotEmpty(t, plainKey)
	assert.Contains(t, plainKey, "sk_")
	assert.Equal(t, "testuser", apiKey.UserID)
	assert.Equal(t, "Test Client", apiKey.Description)
}

func TestSetupTokenGeneration(t *testing.T) {
	// Create auth service
	authService := services.NewAuthService("test-encryption-key-32-chars-long", nil)

	// Generate setup token
	token, err := authService.GenerateSetupToken("testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "testuser", token.UserID)
	assert.Regexp(t, `^[A-Z0-9]{4}-[A-Z0-9]{4}$`, token.Token)
	assert.False(t, token.Used)
	assert.True(t, token.ExpiresAt.After(token.ExpiresAt.Add(-25*time.Hour))) // Expires in ~24 hours
}

func TestSetupTokenValidation(t *testing.T) {
	// Create auth service
	authService := services.NewAuthService("test-encryption-key-32-chars-long", nil)

	// Generate setup token
	token, err := authService.GenerateSetupToken("testuser")
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
