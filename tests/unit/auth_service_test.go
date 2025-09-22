package unit

import (
	"os"
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-secret", store)

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
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-secret", store)

	user, _ := authService.Authenticate("testuser", "testpass123")
	token, err := authService.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	store := storage.NewMemoryStorage()
	authService := services.NewAuthService("test-secret", store)

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

func TestNewEnvironmentConfiguration(t *testing.T) {
	config := models.NewEnvironmentConfiguration()

	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "development", config.Environment)
	assert.Empty(t, config.JWT_SECRET)
}

func TestLoadFromEnv_Valid(t *testing.T) {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long")
	os.Setenv("PORT", "9090")
	os.Setenv("ENVIRONMENT", "production")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
	}()

	config := models.NewEnvironmentConfiguration()
	err := config.LoadFromEnv(os.Getenv)

	assert.NoError(t, err)
	assert.Equal(t, "test-secret-key-32-chars-long", config.JWT_SECRET)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "production", config.Environment)
}

func TestValidate_Valid(t *testing.T) {
	config := &models.EnvironmentConfiguration{
		JWT_SECRET:  "test-secret-key-32-chars-long-enough",
		Port:        8080,
		Environment: "development",
	}

	err := config.Validate()

	assert.NoError(t, err)
}

func TestIsProduction(t *testing.T) {
	config := &models.EnvironmentConfiguration{Environment: "production"}
	assert.True(t, config.IsProduction())

	config.Environment = "development"
	assert.False(t, config.IsProduction())
}
