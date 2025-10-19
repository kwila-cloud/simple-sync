package storage

import (
	"errors"
	"simple-sync/src/models"
	"testing"
)

// Database-specific error types
var (
	ErrNotFound     = errors.New("resource not found")
	ErrDuplicateKey = errors.New("duplicate key")
	ErrInvalidData  = errors.New("invalid data")

	// More specific error types
	ErrUserNotFound         = errors.New("user not found")
	ErrApiKeyNotFound       = errors.New("API key not found")
	ErrSetupTokenNotFound   = errors.New("setup token not found")
	ErrInvalidApiKeyFormat  = errors.New("invalid API key format")
	ErrInvalidApiKey        = errors.New("invalid API key")
	ErrInvalidSetupToken    = errors.New("invalid setup token")
	ErrSetupTokenExpired    = errors.New("setup token is expired or already used")
	ErrInvalidTimestamp     = errors.New("invalid timestamp")
	ErrInvalidUserPattern   = errors.New("invalid user pattern")
	ErrInvalidItemPattern   = errors.New("invalid item pattern")
	ErrInvalidActionPattern = errors.New("invalid action pattern")
	ErrInvalidAclType       = errors.New("type must be either 'allow' or 'deny'")

	// ACL validation specific errors
	ErrAclUserEmpty               = errors.New("user pattern cannot be empty")
	ErrAclItemEmpty               = errors.New("item pattern cannot be empty")
	ErrAclActionEmpty             = errors.New("action pattern cannot be empty")
	ErrAclUserMultipleWildcards   = errors.New("user pattern can have at most one wildcard at the end")
	ErrAclItemMultipleWildcards   = errors.New("item pattern can have at most one wildcard at the end")
	ErrAclActionMultipleWildcards = errors.New("action pattern can have at most one wildcard at the end")
	ErrAclUserControlChars        = errors.New("user pattern contains invalid control characters")
	ErrAclItemControlChars        = errors.New("item pattern contains invalid control characters")
	ErrAclActionControlChars      = errors.New("action pattern contains invalid control characters")
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
