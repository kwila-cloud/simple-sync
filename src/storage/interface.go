package storage

import "simple-sync/src/models"

// Storage defines the interface for data persistence
type Storage interface {
	// Event operations
	SaveEvents(events []models.Event) error
	LoadEvents(fromTimestamp *uint64) ([]models.Event, error)

	// User operations
	SaveUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUUID(uuid string) (*models.User, error)

	// API Key operations
	CreateAPIKey(apiKey *models.APIKey) error
	GetAPIKeyByHash(hash string) (*models.APIKey, error)
	UpdateAPIKey(apiKey *models.APIKey) error

	// Setup Token operations
	CreateSetupToken(token *models.SetupToken) error
	GetSetupToken(token string) (*models.SetupToken, error)
	UpdateSetupToken(token *models.SetupToken) error
	InvalidateUserSetupTokens(userID string) error
}
