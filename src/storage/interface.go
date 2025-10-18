package storage

import "simple-sync/src/models"

// Storage defines the interface for data persistence
type Storage interface {
	// Event operations
	SaveEvents(events []models.Event) error
	LoadEvents() ([]models.Event, error)

	// User operations
	SaveUser(user *models.User) error
	GetUserById(id string) (*models.User, error)

	// API Key operations
	CreateApiKey(apiKey *models.APIKey) error
	GetApiKeyByHash(hash string) (*models.APIKey, error)
	GetAllApiKeys() ([]*models.APIKey, error)
	UpdateApiKey(apiKey *models.APIKey) error
	InvalidateUserApiKeys(userID string) error

	// Setup Token operations
	CreateSetupToken(token *models.SetupToken) error
	GetSetupToken(token string) (*models.SetupToken, error)
	UpdateSetupToken(token *models.SetupToken) error
	InvalidateUserSetupTokens(userID string) error

	// ACL operations
	CreateAclRule(rule *models.AclRule) error
	GetAclRules() ([]models.AclRule, error)
}
