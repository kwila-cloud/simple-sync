package unit

import (
	"strings"
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/storage"
)

func TestSQLiteCreateAclRuleAndGetAclRules(t *testing.T) {
	store := storage.NewSQLiteStorage()
	if err := store.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize sqlite storage: %v", err)
	}
	defer store.Close()

	rule := models.AclRule{
		User:   "sqlite-user",
		Item:   "sqlite-item",
		Action: "read",
		Type:   "allow",
	}

	// Create rule
	if err := store.CreateAclRule(&rule); err != nil {
		t.Fatalf("expected no error creating rule, got %v", err)
	}

	// Retrieve
	rules, err := store.GetAclRules()
	if err != nil {
		t.Fatalf("expected no error getting rules, got %v", err)
	}

	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	if rules[0].User != rule.User {
		t.Errorf("expected user %s, got %s", rule.User, rules[0].User)
	}
	if rules[0].Item != rule.Item {
		t.Errorf("expected item %s, got %s", rule.Item, rules[0].Item)
	}
	if rules[0].Action != rule.Action {
		t.Errorf("expected action %s, got %s", rule.Action, rules[0].Action)
	}
	if rules[0].Type != rule.Type {
		t.Errorf("expected type %s, got %s", rule.Type, rules[0].Type)
	}
}

func TestSQLiteCreateAclRuleDuplicate(t *testing.T) {
	store := storage.NewSQLiteStorage()
	if err := store.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize sqlite storage: %v", err)
	}
	defer store.Close()

	rule := models.AclRule{
		User:   "dup-user",
		Item:   "dup-item",
		Action: "write",
		Type:   "allow",
	}

	if err := store.CreateAclRule(&rule); err != nil {
		t.Fatalf("expected no error creating rule first time, got %v", err)
	}

	// Second insert should result in duplicate key error
	if err := store.CreateAclRule(&rule); err == nil {
		t.Fatalf("expected duplicate key error, got nil")
	} else {
		if !strings.Contains(err.Error(), "duplicate") && err != storage.ErrDuplicateKey {
			t.Errorf("expected ErrDuplicateKey or message containing 'duplicate', got %v", err)
		}
	}
}

func TestSQLiteGetAclRulesEmpty(t *testing.T) {
	store := storage.NewSQLiteStorage()
	if err := store.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize sqlite storage: %v", err)
	}
	defer store.Close()

	rules, err := store.GetAclRules()
	if err != nil {
		t.Fatalf("expected no error getting rules from empty db, got %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(rules))
	}
}
