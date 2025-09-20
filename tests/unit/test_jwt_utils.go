package unit

import (
	"testing"

	"simple-sync/src/services"
	"simple-sync/src/storage"
	"simple-sync/src/utils"

	"github.com/stretchr/testify/assert"
)

func TestExtractTokenFromHeader(t *testing.T) {
	// Test valid Bearer token
	token, err := utils.ExtractTokenFromHeader("Bearer valid-token")
	assert.NoError(t, err)
	assert.Equal(t, "valid-token", token)

	// Test missing Bearer
	_, err = utils.ExtractTokenFromHeader("valid-token")
	assert.Error(t, err)

	// Test invalid format
	_, err = utils.ExtractTokenFromHeader("Bearer")
	assert.Error(t, err)

	// Test empty header
	_, err = utils.ExtractTokenFromHeader("")
	assert.Error(t, err)
}

func TestValidateAndExtractClaims(t *testing.T) {
	// Setup storage and auth service
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-secret", store)

	// Get a valid token
	user, _ := authService.Authenticate("testuser", "testpass123")
	token, _ := authService.GenerateToken(user)

	// Test valid token
	claims, err := utils.ValidateAndExtractClaims(token, authService)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserUUID)
	assert.Equal(t, "testuser", claims.Username)

	// Test invalid token
	_, err = utils.ValidateAndExtractClaims("invalid-token", authService)
	assert.Error(t, err)
}
