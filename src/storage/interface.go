package storage

import (
	"errors"
	"simple-sync/src/models"
	"testing"
)

// Storage-specific error types
var (
	ErrNotFound           = errors.New("resource not found")
	ErrDuplicateKey       = errors.New("duplicate key")
	ErrInvalidData        = errors.New("invalid data")
	ErrApiKeyNotFound     = errors.New("API key not found")
	ErrSetupTokenNotFound = errors.New("setup token not found")
)

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

// NewStorage creates a new storage instance based on the current environment
// Returns TestStorage when running tests, SQLiteStorage in production (future)
func NewStorage() Storage {
	if testing.Testing() {
		return NewTestStorage(nil)
	}
	// TODO: Return SQLiteStorage for production
	return NewTestStorage(nil)
}

// NewStorageWithAclRules creates a new storage instance with initial ACL rules
// Returns TestStorage when running tests, SQLiteStorage in production (future)
func NewStorageWithAclRules(aclRules []models.AclRule) Storage {
	if testing.Testing() {
		return NewTestStorage(aclRules)
	}
	// TODO: Return SQLiteStorage for production
	return NewTestStorage(aclRules)
}
