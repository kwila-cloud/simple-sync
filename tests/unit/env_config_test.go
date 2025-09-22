package unit

import (
	"os"
	"testing"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

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

func TestValidate_Port80Allowed(t *testing.T) {
	config := &models.EnvironmentConfiguration{
		JWT_SECRET:  "test-secret-key-32-chars-long-enough",
		Port:        80,
		Environment: "development",
	}

	err := config.Validate()

	assert.NoError(t, err)
}

func TestLoadFromEnv_PortTooLow(t *testing.T) {
	// Set up test environment with port below 80
	os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long")
	os.Setenv("PORT", "79")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
	}()

	config := models.NewEnvironmentConfiguration()
	err := config.LoadFromEnv(os.Getenv)

	if err == nil {
		t.Error("Expected error for port too low")
	}

	if err.Error() != "PORT must be between 80 and 65535" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestValidate_PortTooLow(t *testing.T) {
	config := &models.EnvironmentConfiguration{
		JWT_SECRET:  "test-secret-key-32-chars-long-enough",
		Port:        79,
		Environment: "development",
	}

	err := config.Validate()

	if err == nil {
		t.Error("Expected error for port too low")
	}

	if err.Error() != "PORT must be between 80 and 65535" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}
