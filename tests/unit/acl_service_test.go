package unit

import (
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestAclService_CheckPermission(t *testing.T) {
	store := storage.NewMemoryStorage()
	aclService := services.NewAclService(store)

	// Test root bypass
	assert.True(t, aclService.CheckPermission(".root", "any", "any"))

	// Test deny by default
	assert.False(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Add allow rule
	rule := models.AclRule{
		User:      "user1",
		Item:      "item1",
		Action:    "action1",
		Type:      "allow",
		Timestamp: 1000,
	}
	aclService.AddRule(rule)

	// Now should allow
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Different user should deny
	assert.False(t, aclService.CheckPermission("user2", "item1", "action1"))
}

func TestAclService_Matches(t *testing.T) {
	store := storage.NewMemoryStorage()
	aclService := services.NewAclService(store)

	// Test deny by default when no rules
	assert.False(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Test matches function indirectly via permission with wildcard
	rule := models.AclRule{
		User:      "*",
		Item:      "item1",
		Action:    "action1",
		Type:      "allow",
		Timestamp: 1000,
	}
	aclService.AddRule(rule)

	assert.True(t, aclService.CheckPermission("anyuser", "item1", "action1"))
}

func TestAclService_Specificity(t *testing.T) {
	store := storage.NewMemoryStorage()
	aclService := services.NewAclService(store)

	// Add deny rule with lower specificity
	denyRule := models.AclRule{
		User:      "*",
		Item:      "*",
		Action:    "*",
		Type:      "deny",
		Timestamp: 1000,
	}
	aclService.AddRule(denyRule)

	// Add allow rule with higher specificity
	allowRule := models.AclRule{
		User:      "user1",
		Item:      "item1",
		Action:    "action1",
		Type:      "allow",
		Timestamp: 1001,
	}
	aclService.AddRule(allowRule)

	// Should allow due to higher specificity
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
}

func TestAclService_TimestampResolution(t *testing.T) {
	store := storage.NewMemoryStorage()
	aclService := services.NewAclService(store)

	// Add allow rule
	allowRule := models.AclRule{
		User:      "user1",
		Item:      "item1",
		Action:    "action1",
		Type:      "allow",
		Timestamp: 1000,
	}
	aclService.AddRule(allowRule)

	// Add deny rule with same specificity but later timestamp
	denyRule := models.AclRule{
		User:      "user1",
		Item:      "item1",
		Action:    "action1",
		Type:      "deny",
		Timestamp: 1001,
	}
	aclService.AddRule(denyRule)

	// Should deny due to later timestamp
	assert.False(t, aclService.CheckPermission("user1", "item1", "action1"))
}

func TestAclService_ValidateAclEvent(t *testing.T) {
	store := storage.NewMemoryStorage()
	aclService := services.NewAclService(store)

	// Test non-ACL event (should return false)
	nonAclEvent := &models.Event{
		Item:   "some-item",
		Action: "some-action",
		User:   "user1",
	}
	assert.False(t, aclService.ValidateAclEvent(nonAclEvent))

	// Test invalid ACL action (not .acl.allow or .acl.deny)
	invalidActionEvent := &models.Event{
		Item:    ".acl",
		Action:  ".acl.invalid",
		User:    "user1",
		Payload: `{"user":"user2","item":"item1","action":"read","type":"allow"}`,
	}
	assert.False(t, aclService.ValidateAclEvent(invalidActionEvent))

	// Test malformed JSON payload
	malformedEvent := &models.Event{
		Item:    ".acl",
		Action:  ".acl.allow",
		User:    "user1",
		Payload: `{"user":"user2","item":"item1","action":"read","type":"allow"`, // missing closing brace
	}
	assert.False(t, aclService.ValidateAclEvent(malformedEvent))

	// Test empty fields in ACL rule
	emptyFieldsEvent := &models.Event{
		Item:    ".acl",
		Action:  ".acl.allow",
		User:    "user1",
		Payload: `{"user":"","item":"item1","action":"read","type":"allow"}`,
	}
	assert.False(t, aclService.ValidateAclEvent(emptyFieldsEvent))

	// Test invalid wildcard patterns (multiple wildcards)
	invalidWildcardEvent := &models.Event{
		Item:    ".acl",
		Action:  ".acl.allow",
		User:    "user1",
		Payload: `{"user":"user*test*","item":"item1","action":"read","type":"allow"}`,
	}
	assert.False(t, aclService.ValidateAclEvent(invalidWildcardEvent))

	// Test invalid wildcard patterns (wildcard not at end)
	invalidWildcardEvent2 := &models.Event{
		Item:    ".acl",
		Action:  ".acl.allow",
		User:    "user1",
		Payload: `{"user":"user*test","item":"item1","action":"read","type":"allow"}`,
	}
	assert.False(t, aclService.ValidateAclEvent(invalidWildcardEvent2))

	// Test permission check - user without permission
	noPermissionEvent := &models.Event{
		Item:    ".acl",
		Action:  ".acl.allow",
		User:    "user1",
		Payload: `{"user":"user2","item":"item1","action":"read","type":"allow"}`,
	}
	assert.False(t, aclService.ValidateAclEvent(noPermissionEvent))

	// Test valid ACL event with root user (should bypass permission check)
	rootEvent := &models.Event{
		Item:    ".acl",
		Action:  ".acl.allow",
		User:    ".root",
		Payload: `{"user":"user2","item":"item1","action":"read","type":"allow"}`,
	}
	assert.True(t, aclService.ValidateAclEvent(rootEvent))

	// Test valid ACL event with proper permissions
	// First, add permission for user1 to set .acl.allow
	aclRule := models.AclRule{
		User:      "user1",
		Item:      ".acl",
		Action:    ".acl.allow",
		Type:      "allow",
		Timestamp: 1000,
	}
	aclService.AddRule(aclRule)

	validEvent := &models.Event{
		Item:    ".acl",
		Action:  ".acl.allow",
		User:    "user1",
		Payload: `{"user":"user2","item":"item1","action":"read","type":"allow"}`,
	}
	assert.True(t, aclService.ValidateAclEvent(validEvent))
}
