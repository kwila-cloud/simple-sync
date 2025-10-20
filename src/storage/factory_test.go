package storage

import (
	"simple-sync/src/models"
	"testing"
)

func TestNewStorage(t *testing.T) {
	// Test that factory returns TestStorage when running tests
	store := NewStorage()
	if store == nil {
		t.Fatal("Expected storage to be created")
	}

	// Verify it's TestStorage by checking type
	_, isTestStorage := store.(*TestStorage)
	if !isTestStorage {
		t.Errorf("Expected TestStorage, got %T", store)
	}
}

func TestNewStorageWithAclRules(t *testing.T) {
	// Ensure NewStorage returns TestStorage in test environments
	store := NewStorage()
	if store == nil {
		t.Fatal("Expected storage to be created")
	}

	// Verify it's TestStorage
	_, isTestStorage := store.(*TestStorage)
	if !isTestStorage {
		t.Errorf("Expected TestStorage, got %T", store)
	}

	// Verify ACL seeding can be done explicitly via CreateAclRule
	rule := models.AclRule{User: "user1", Item: "item1", Action: "read", Type: "allow"}
	if err := store.CreateAclRule(&rule); err != nil {
		t.Fatalf("Failed to create ACL rule: %v", err)
	}

	rules, err := store.GetAclRules()
	if err != nil {
		t.Fatalf("Failed to get ACL rules: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("Expected 1 ACL rule, got %d", len(rules))
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
			err:      ErrNotFound,
			expected: "resource not found",
		},
		{
			name:     "ErrDuplicateKey",
			err:      ErrDuplicateKey,
			expected: "duplicate key",
		},
		{
			name:     "ErrInvalidData",
			err:      ErrInvalidData,
			expected: "invalid data",
		},
		{
			name:     "ErrApiKeyNotFound",
			err:      ErrApiKeyNotFound,
			expected: "API key not found",
		},
		{
			name:     "ErrSetupTokenNotFound",
			err:      ErrSetupTokenNotFound,
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
