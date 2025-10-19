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

	// Test error messages
	expectedNotFound := "resource not found"
	if ErrNotFound.Error() != expectedNotFound {
		t.Errorf("Expected error message '%s', got '%s'", expectedNotFound, ErrNotFound.Error())
	}

	expectedDuplicate := "duplicate key"
	if ErrDuplicateKey.Error() != expectedDuplicate {
		t.Errorf("Expected error message '%s', got '%s'", expectedDuplicate, ErrDuplicateKey.Error())
	}

	expectedInvalid := "invalid data"
	if ErrInvalidData.Error() != expectedInvalid {
		t.Errorf("Expected error message '%s', got '%s'", expectedInvalid, ErrInvalidData.Error())
	}
}
