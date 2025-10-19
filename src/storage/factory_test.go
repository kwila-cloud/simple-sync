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
	// Test factory with ACL rules
	aclRules := []models.AclRule{
		{
			User:   "user1",
			Item:   "item1",
			Action: "read",
			Type:   "allow",
		},
	}

	store := NewStorageWithAclRules(aclRules)
	if store == nil {
		t.Fatal("Expected storage to be created")
	}

	// Verify it's TestStorage
	_, isTestStorage := store.(*TestStorage)
	if !isTestStorage {
		t.Errorf("Expected TestStorage, got %T", store)
	}

	// Verify ACL rules were loaded
	rules, err := store.GetAclRules()
	if err != nil {
		t.Fatalf("Failed to get ACL rules: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("Expected 1 ACL rule, got %d", len(rules))
	}
}

func TestErrorTypes(t *testing.T) {
	// Test that error types are properly defined
	if ErrNotFound == nil {
		t.Error("ErrNotFound should not be nil")
	}
	if ErrDuplicateKey == nil {
		t.Error("ErrDuplicateKey should not be nil")
	}
	if ErrInvalidData == nil {
		t.Error("ErrInvalidData should not be nil")
	}

	// Test specific error types
	if ErrUserNotFound == nil {
		t.Error("ErrUserNotFound should not be nil")
	}
	if ErrApiKeyNotFound == nil {
		t.Error("ErrApiKeyNotFound should not be nil")
	}
	if ErrSetupTokenNotFound == nil {
		t.Error("ErrSetupTokenNotFound should not be nil")
	}
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
	if ErrInvalidTimestamp == nil {
		t.Error("ErrInvalidTimestamp should not be nil")
	}
	if ErrInvalidUserPattern == nil {
		t.Error("ErrInvalidUserPattern should not be nil")
	}
	if ErrInvalidItemPattern == nil {
		t.Error("ErrInvalidItemPattern should not be nil")
	}
	if ErrInvalidActionPattern == nil {
		t.Error("ErrInvalidActionPattern should not be nil")
	}
	if ErrInvalidAclType == nil {
		t.Error("ErrInvalidAclType should not be nil")
	}

	// ACL validation specific errors
	if ErrAclUserEmpty == nil {
		t.Error("ErrAclUserEmpty should not be nil")
	}
	if ErrAclItemEmpty == nil {
		t.Error("ErrAclItemEmpty should not be nil")
	}
	if ErrAclActionEmpty == nil {
		t.Error("ErrAclActionEmpty should not be nil")
	}
	if ErrAclUserMultipleWildcards == nil {
		t.Error("ErrAclUserMultipleWildcards should not be nil")
	}
	if ErrAclItemMultipleWildcards == nil {
		t.Error("ErrAclItemMultipleWildcards should not be nil")
	}
	if ErrAclActionMultipleWildcards == nil {
		t.Error("ErrAclActionMultipleWildcards should not be nil")
	}
	if ErrAclUserControlChars == nil {
		t.Error("ErrAclUserControlChars should not be nil")
	}
	if ErrAclItemControlChars == nil {
		t.Error("ErrAclItemControlChars should not be nil")
	}
	if ErrAclActionControlChars == nil {
		t.Error("ErrAclActionControlChars should not be nil")
	}

	// Test error messages
	expectedUserNotFound := "user not found"
	if ErrUserNotFound.Error() != expectedUserNotFound {
		t.Errorf("Expected error message '%s', got '%s'", expectedUserNotFound, ErrUserNotFound.Error())
	}

	expectedApiKeyNotFound := "API key not found"
	if ErrApiKeyNotFound.Error() != expectedApiKeyNotFound {
		t.Errorf("Expected error message '%s', got '%s'", expectedApiKeyNotFound, ErrApiKeyNotFound.Error())
	}

	expectedInvalidTimestamp := "invalid timestamp"
	if ErrInvalidTimestamp.Error() != expectedInvalidTimestamp {
		t.Errorf("Expected error message '%s', got '%s'", expectedInvalidTimestamp, ErrInvalidTimestamp.Error())
	}
}
