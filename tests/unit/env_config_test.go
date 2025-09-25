package unit

import (
	"testing"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

// testEnv provides an abstraction for environment variables in tests
type testEnv struct {
	vars map[string]string
}

func newTestEnv() *testEnv {
	return &testEnv{
		vars: make(map[string]string),
	}
}

func (te *testEnv) set(key, value string) {
	te.vars[key] = value
}

func (te *testEnv) unset(key string) {
	delete(te.vars, key)
}

func (te *testEnv) get(key string) string {
	return te.vars[key]
}

func TestNewEnvironmentConfiguration(t *testing.T) {
	config := models.NewEnvironmentConfiguration()

	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "development", config.Environment)
}

func TestLoadFromEnv_Valid(t *testing.T) {
	// Set up test environment
	env := newTestEnv()
	env.set("PORT", "9090")
	env.set("ENVIRONMENT", "production")

	config := models.NewEnvironmentConfiguration()
	err := config.LoadFromEnv(env.get)

	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "production", config.Environment)
}

func TestValidate_Valid(t *testing.T) {
	config := &models.EnvironmentConfiguration{
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
		Port:        80,
		Environment: "development",
	}

	err := config.Validate()

	assert.NoError(t, err)
}

func TestLoadFromEnv_PortTooLow(t *testing.T) {
	// Set up test environment with port below 80
	env := newTestEnv()
	env.set("PORT", "79")

	config := models.NewEnvironmentConfiguration()
	err := config.LoadFromEnv(env.get)

	if err == nil {
		t.Error("Expected error for port too low")
	}

	if err.Error() != "PORT must be between 80 and 65535" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestValidate_PortTooLow(t *testing.T) {
	config := &models.EnvironmentConfiguration{
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
