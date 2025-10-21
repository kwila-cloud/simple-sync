package unit

import (
	"simple-sync/src/errors"
	"testing"
)

func TestServiceErrorTypes(t *testing.T) {
	// Test that service error types are properly defined
	if errors.ErrInvalidApiKeyFormat == nil {
		t.Error("ErrInvalidApiKeyFormat should not be nil")
	}
	if errors.ErrInvalidApiKey == nil {
		t.Error("ErrInvalidApiKey should not be nil")
	}
	if errors.ErrInvalidSetupToken == nil {
		t.Error("ErrInvalidSetupToken should not be nil")
	}
	if errors.ErrSetupTokenExpired == nil {
		t.Error("ErrSetupTokenExpired should not be nil")
	}
	if errors.ErrUserNotFound == nil {
		t.Error("ErrUserNotFound should not be nil")
	}

	// Test error messages
	expectedUserNotFound := "user not found"
	if errors.ErrUserNotFound.Error() != expectedUserNotFound {
		t.Errorf("Expected error message '%s', got '%s'", expectedUserNotFound, errors.ErrUserNotFound.Error())
	}

	expectedInvalidApiKey := "invalid API key"
	if errors.ErrInvalidApiKey.Error() != expectedInvalidApiKey {
		t.Errorf("Expected error message '%s', got '%s'", expectedInvalidApiKey, errors.ErrInvalidApiKey.Error())
	}
}
