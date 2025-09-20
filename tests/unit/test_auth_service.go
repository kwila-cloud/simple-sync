package unit

import (
	"testing"

	"simple-sync/src/services"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	authService := services.NewAuthService("test-secret")

	// Test valid credentials
	user, err := authService.Authenticate("testuser", "testpass123")
	assert.NoError(t, err)
	assert.Equal(t, "user-123", user.UUID)
	assert.Equal(t, "testuser", user.Username)

	// Test invalid username
	_, err = authService.Authenticate("wronguser", "testpass123")
	assert.Error(t, err)

	// Test invalid password
	_, err = authService.Authenticate("testuser", "wrongpass")
	assert.Error(t, err)
}

func TestGenerateToken(t *testing.T) {
	authService := services.NewAuthService("test-secret")

	user, _ := authService.Authenticate("testuser", "testpass123")
	token, err := authService.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	authService := services.NewAuthService("test-secret")

	user, _ := authService.Authenticate("testuser", "testpass123")
	token, _ := authService.GenerateToken(user)

	// Test valid token
	claims, err := authService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserUUID)
	assert.Equal(t, "testuser", claims.Username)

	// Test invalid token
	_, err = authService.ValidateToken("invalid-token")
	assert.Error(t, err)

	// Test expired token (simulate by creating old token)
	// Note: For full test, would need to mock time
}
