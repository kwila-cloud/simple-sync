package unit

import (
	"strings"
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/storage"
)

func TestCreateAclRule(t *testing.T) {
	testStorage := storage.NewTestStorage([]models.AclRule{})

	rule := models.AclRule{
		User:   "test-user",
		Item:   "test-item",
		Action: "read",
		Type:   "allow",
	}

	err := testStorage.AddAclRule(&rule)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the rule was stored by checking events
	events, err := testStorage.LoadEvents()
	if err != nil {
		t.Fatalf("Expected no error loading events, got %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if !events[0].IsAclEvent() {
		t.Fatalf("Expected ACL event, got regular event")
	}

	storedRule, err := events[0].ToAclRule()
	if err != nil {
		t.Fatalf("Expected no error converting to ACL rule, got %v", err)
	}

	if storedRule.User != rule.User {
		t.Errorf("Expected user %s, got %s", rule.User, storedRule.User)
	}

	if storedRule.Item != rule.Item {
		t.Errorf("Expected item %s, got %s", rule.Item, storedRule.Item)
	}

	if storedRule.Action != rule.Action {
		t.Errorf("Expected action %s, got %s", rule.Action, storedRule.Action)
	}

	if storedRule.Type != rule.Type {
		t.Errorf("Expected type %s, got %s", rule.Type, storedRule.Type)
	}
}

func TestGetAclRules(t *testing.T) {
	initialRules := []models.AclRule{
		{
			User:   "user1",
			Item:   "item1",
			Action: "read",
			Type:   "allow",
		},
		{
			User:   "user2",
			Item:   "item2",
			Action: "write",
			Type:   "deny",
		},
	}

	testStorage := storage.NewTestStorage(initialRules)

	rules, err := testStorage.GetAclRules()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(rules) != 2 {
		t.Fatalf("Expected 2 rules, got %d", len(rules))
	}

	// Verify first rule
	if rules[0].User != "user1" {
		t.Errorf("Expected first rule user user1, got %s", rules[0].User)
	}

	if rules[0].Item != "item1" {
		t.Errorf("Expected first rule item item1, got %s", rules[0].Item)
	}

	// Verify second rule
	if rules[1].User != "user2" {
		t.Errorf("Expected second rule user user2, got %s", rules[1].User)
	}

	if rules[1].Action != "write" {
		t.Errorf("Expected second rule action write, got %s", rules[1].Action)
	}
}

func TestGetAclRulesEmpty(t *testing.T) {
	testStorage := storage.NewTestStorage([]models.AclRule{})

	rules, err := testStorage.GetAclRules()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(rules) != 0 {
		t.Fatalf("Expected 0 rules, got %d", len(rules))
	}
}

func TestCreateAclRuleAndGetAclRulesIntegration(t *testing.T) {
	testStorage := storage.NewTestStorage([]models.AclRule{})

	// Add multiple rules
	rules := []models.AclRule{
		{
			User:   "alice",
			Item:   "document1",
			Action: "read",
			Type:   "allow",
		},
		{
			User:   "bob",
			Item:   "document2",
			Action: "write",
			Type:   "allow",
		},
		{
			User:   "charlie",
			Item:   "document3",
			Action: "delete",
			Type:   "deny",
		},
	}

	for i := range rules {
		err := testStorage.AddAclRule(&rules[i])
		if err != nil {
			t.Fatalf("Expected no error creating rule %d, got %v", i, err)
		}
	}

	// Retrieve all rules
	storedRules, err := testStorage.GetAclRules()
	if err != nil {
		t.Fatalf("Expected no error getting rules, got %v", err)
	}

	if len(storedRules) != 3 {
		t.Fatalf("Expected 3 rules, got %d", len(storedRules))
	}

	// Verify all rules were stored correctly
	for i, expectedRule := range rules {
		if storedRules[i].User != expectedRule.User {
			t.Errorf("Rule %d: expected user %s, got %s", i, expectedRule.User, storedRules[i].User)
		}

		if storedRules[i].Item != expectedRule.Item {
			t.Errorf("Rule %d: expected item %s, got %s", i, expectedRule.Item, storedRules[i].Item)
		}

		if storedRules[i].Action != expectedRule.Action {
			t.Errorf("Rule %d: expected action %s, got %s", i, expectedRule.Action, storedRules[i].Action)
		}

		if storedRules[i].Type != expectedRule.Type {
			t.Errorf("Rule %d: expected type %s, got %s", i, expectedRule.Type, storedRules[i].Type)
		}
	}
}

func TestGetAclRulesWithMalformedRule(t *testing.T) {
	testStorage := storage.NewTestStorage([]models.AclRule{})

	// Add a valid rule
	validRule := models.AclRule{
		User:   "alice",
		Item:   "document1",
		Action: "read",
		Type:   "allow",
	}
	err := testStorage.CreateAclRule(&validRule)
	if err != nil {
		t.Fatalf("Expected no error creating valid rule, got %v", err)
	}

	// Manually add a malformed ACL event (invalid JSON)
	malformedEvent := models.Event{
		User:    ".root",
		Item:    ".acl",
		Action:  ".acl.addRule",
		Payload: "{invalid json}",
	}
	err = testStorage.AddEvents([]models.Event{malformedEvent})
	if err != nil {
		t.Fatalf("Expected no error saving malformed event, got %v", err)
	}

	// GetAclRules should return an error due to malformed rule
	_, err = testStorage.GetAclRules()
	if err == nil {
		t.Fatalf("Expected error due to malformed ACL rule, got nil")
	}

	expectedError := "malformed ACL rule in event"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}
