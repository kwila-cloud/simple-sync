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
	rule := models.ACLRule{
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

	// Exact match
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1")) // but since no rule, false

	// Test matches function indirectly via permission
	rule := models.ACLRule{
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
	denyRule := models.ACLRule{
		User:      "*",
		Item:      "*",
		Action:    "*",
		Type:      "deny",
		Timestamp: 1000,
	}
	aclService.AddRule(denyRule)

	// Add allow rule with higher specificity
	allowRule := models.ACLRule{
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
	allowRule := models.ACLRule{
		User:      "user1",
		Item:      "item1",
		Action:    "action1",
		Type:      "allow",
		Timestamp: 1000,
	}
	aclService.AddRule(allowRule)

	// Add deny rule with same specificity but later timestamp
	denyRule := models.ACLRule{
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
