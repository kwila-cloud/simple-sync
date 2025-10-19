package services

import (
	"simple-sync/src/storage"
	"testing"
)

func TestServiceErrorTypes(t *testing.T) {
	// Test that service error types are properly defined
	if ErrInvalidApiKeyFormat == nil {
		t.Error("ErrInvalidApiKeyFormat should not be nil")
	}
	if ErrInvalidApiKey == nil {
		t.Error("ErrInvalidApiKey should not be nil")
	}
	if ErrInvalidSetupToken == nil {
		t.Error("ErrInvalidSetupToken should not be nil")
	}
	if ErrSetupTokenExpired == nil {
		t.Error("ErrSetupTokenExpired should not be nil")
	}
	if ErrUserNotFound == nil {
		t.Error("ErrUserNotFound should not be nil")
	}

	// Test error messages
	expectedUserNotFound := "user not found"
	if ErrUserNotFound.Error() != expectedUserNotFound {
		t.Errorf("Expected error message '%s', got '%s'", expectedUserNotFound, ErrUserNotFound.Error())
	}

	expectedInvalidApiKey := "invalid API key"
	if ErrInvalidApiKey.Error() != expectedInvalidApiKey {
		t.Errorf("Expected error message '%s', got '%s'", expectedInvalidApiKey, ErrInvalidApiKey.Error())
	}
}

func TestErrorTranslation(t *testing.T) {
	// Test that storage.ErrNotFound gets translated to ErrUserNotFound
	if storage.ErrNotFound == nil {
		t.Error("storage.ErrNotFound should not be nil")
	}

	// The translation logic is tested in the auth service integration tests
	// This just verifies the error types exist and have correct messages
}
