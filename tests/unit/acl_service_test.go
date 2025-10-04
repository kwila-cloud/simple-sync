package unit

import (
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestAclService_CheckPermission(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService := services.NewAclService(store)

	// Test root bypass
	assert.True(t, aclService.CheckPermission(".root", "any", "any"))

	// Test deny by default
	assert.False(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Add allow rule
	rule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	aclService.AddRule(rule)

	// Now should allow
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Different user should deny
	assert.False(t, aclService.CheckPermission("user2", "item1", "action1"))
}

func TestAclService_Matches(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService := services.NewAclService(store)

	// Test deny by default when no rules
	assert.False(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Test matches function indirectly via permission with wildcard
	rule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	aclService.AddRule(rule)

	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
}

func TestAclService_Specificity(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService := services.NewAclService(store)

	// Add deny rule with lower specificity
	denyRule := models.AclRule{
		User:   "user1",
		Item:   "*",
		Action: "*",
		Type:   "deny",
	}
	aclService.AddRule(denyRule)

	// Add allow rule with higher specificity
	allowRule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	aclService.AddRule(allowRule)

	// Should allow due to higher specificity
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
}

func TestAclService_OrderResolution(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService := services.NewAclService(store)

	// Add deny rule
	allowRule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "deny",
	}
	aclService.AddRule(allowRule)

	// Add allow rule with same specificity
	denyRule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	aclService.AddRule(denyRule)

	// Should allow due to later event
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
}
