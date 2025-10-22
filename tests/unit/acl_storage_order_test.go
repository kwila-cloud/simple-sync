package unit

import (
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/storage"
)

// Test that GetAclRules returns rules in insertion order
func TestSQLiteAclRulesPreserveInsertionOrder(t *testing.T) {
	store := storage.NewSQLiteStorage()
	if err := store.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize sqlite storage: %v", err)
	}
	defer store.Close()

	rulesToInsert := []models.AclRule{
		{User: "u1", Item: "i1", Action: "a1", Type: "allow"},
		{User: "u2", Item: "i2", Action: "a2", Type: "deny"},
		{User: "u3", Item: "i3", Action: "a3", Type: "allow"},
	}

	for i := range rulesToInsert {
		if err := store.CreateAclRule(&rulesToInsert[i]); err != nil {
			t.Fatalf("failed to create rule %d: %v", i, err)
		}
	}

	got, err := store.GetAclRules()
	if err != nil {
		t.Fatalf("expected no error getting rules, got %v", err)
	}

	if len(got) != len(rulesToInsert) {
		t.Fatalf("expected %d rules, got %d", len(rulesToInsert), len(got))
	}

	for i := range rulesToInsert {
		expected := rulesToInsert[i]
		actual := got[i]
		if actual.User != expected.User || actual.Item != expected.Item || actual.Action != expected.Action || actual.Type != expected.Type {
			t.Fatalf("rule %d mismatch: expected %+v, got %+v", i, expected, actual)
		}
	}
}
