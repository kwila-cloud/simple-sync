package unit

import (
	"simple-sync/src/storage"
	"testing"
)

func TestNewStorage(t *testing.T) {
	// Test that factory returns TestStorage when running tests
	store := storage.NewStorage()
	if store == nil {
		t.Fatal("Expected storage to be created")
	}

	// Verify it's TestStorage by checking type
	_, isTestStorage := store.(*storage.TestStorage)
	if !isTestStorage {
		t.Errorf("Expected TestStorage, got %T", store)
	}
}

func TestErrorTypes(t *testing.T) {
	// Test error messages for storage-specific errors
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNotFound",
			err:      storage.ErrNotFound,
			expected: "resource not found",
		},
		{
			name:     "ErrDuplicateKey",
			err:      storage.ErrDuplicateKey,
			expected: "duplicate key",
		},
		{
			name:     "ErrInvalidData",
			err:      storage.ErrInvalidData,
			expected: "invalid data",
		},
		{
			name:     "ErrApiKeyNotFound",
			err:      storage.ErrApiKeyNotFound,
			expected: "API key not found",
		},
		{
			name:     "ErrSetupTokenNotFound",
			err:      storage.ErrSetupTokenNotFound,
			expected: "setup token not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, tt.err.Error())
			}
		})
	}
}
