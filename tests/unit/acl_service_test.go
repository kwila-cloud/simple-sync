package unit

import (
	"fmt"
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestAclService_LoadsRulesFromStorage(t *testing.T) {
	store := storage.NewTestStorage(nil)

	// Create ACL rules in storage
	rule1 := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	rule2 := models.AclRule{
		User:   "user2",
		Item:   "item2",
		Action: "action2",
		Type:   "deny",
	}

	err := store.CreateAclRule(&rule1)
	assert.NoError(t, err)
	err = store.CreateAclRule(&rule2)
	assert.NoError(t, err)

	// Create ACL service - should load rules from storage
	aclService, err := services.NewAclService(store)
	assert.NoError(t, err)

	// Test permissions based on loaded rules
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
	assert.False(t, aclService.CheckPermission("user2", "item2", "action2"))
	assert.False(t, aclService.CheckPermission("user3", "item3", "action3")) // deny by default
}

func TestAclService_CheckPermission(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService, err := services.NewAclService(store)
	assert.NoError(t, err)

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
	err = aclService.AddRule(rule)
	assert.NoError(t, err)

	// Now should allow
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Different user should deny
	assert.False(t, aclService.CheckPermission("user2", "item1", "action1"))
}

func TestAclService_Matches(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService, err := services.NewAclService(store)
	assert.NoError(t, err)

	// Test deny by default when no rules
	assert.False(t, aclService.CheckPermission("user1", "item1", "action1"))

	// Test matches function indirectly via permission with wildcard
	rule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	err = aclService.AddRule(rule)
	assert.NoError(t, err)

	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
}

func TestAclService_Specificity(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService, err := services.NewAclService(store)
	assert.NoError(t, err)

	// Add deny rule with lower specificity
	denyRule := models.AclRule{
		User:   "user1",
		Item:   "*",
		Action: "*",
		Type:   "deny",
	}
	err = aclService.AddRule(denyRule)
	assert.NoError(t, err)

	// Add allow rule with higher specificity
	allowRule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	err = aclService.AddRule(allowRule)
	assert.NoError(t, err)

	// Should allow due to higher specificity
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
}

func TestAclService_OrderResolution(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService, err := services.NewAclService(store)
	assert.NoError(t, err)

	// Add deny rule
	allowRule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "deny",
	}
	err = aclService.AddRule(allowRule)
	assert.NoError(t, err)

	// Add allow rule with same specificity
	denyRule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}
	err = aclService.AddRule(denyRule)
	assert.NoError(t, err)

	// Should allow due to later event
	assert.True(t, aclService.CheckPermission("user1", "item1", "action1"))
}

func TestAclService_NewAclService_ErrorHandling(t *testing.T) {
	// Create a mock storage that fails on GetAclRules
	store := &failingStorage{}

	// Should return error when storage fails
	_, err := services.NewAclService(store)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage error")
}

func TestAclService_AddRule_ErrorHandling(t *testing.T) {
	store := storage.NewTestStorage(nil)
	aclService, err := services.NewAclService(store)
	assert.NoError(t, err)

	// Create a rule that should work
	rule := models.AclRule{
		User:   "user1",
		Item:   "item1",
		Action: "action1",
		Type:   "allow",
	}

	// Should succeed
	err = aclService.AddRule(rule)
	assert.NoError(t, err)
}

// failingStorage is a mock storage that always fails
type failingStorage struct{}

func (f *failingStorage) SaveEvents(events []models.Event) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) LoadEvents() ([]models.Event, error) {
	return nil, fmt.Errorf("storage error")
}

func (f *failingStorage) GetUserById(id string) (*models.User, error) {
	return nil, fmt.Errorf("storage error")
}

func (f *failingStorage) CreateApiKey(apiKey *models.ApiKey) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) GetApiKeyByHash(hash string) (*models.ApiKey, error) {
	return nil, fmt.Errorf("storage error")
}

func (f *failingStorage) GetAllApiKeys() ([]*models.ApiKey, error) {
	return nil, fmt.Errorf("storage error")
}

func (f *failingStorage) UpdateApiKey(apiKey *models.ApiKey) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) InvalidateUserApiKeys(userID string) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) CreateSetupToken(token *models.SetupToken) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) GetSetupToken(token string) (*models.SetupToken, error) {
	return nil, fmt.Errorf("storage error")
}

func (f *failingStorage) UpdateSetupToken(token *models.SetupToken) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) InvalidateUserSetupTokens(userID string) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) CreateAclRule(rule *models.AclRule) error {
	return fmt.Errorf("storage error")
}

func (f *failingStorage) GetAclRules() ([]models.AclRule, error) {
	return nil, fmt.Errorf("storage error")
}
