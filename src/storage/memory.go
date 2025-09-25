package storage

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"simple-sync/src/models"
)

// MemoryStorage implements in-memory storage
type MemoryStorage struct {
	events      []models.Event
	users       map[string]*models.User       // id -> user
	apiKeys     map[string]*models.APIKey     // uuid -> api key
	setupTokens map[string]*models.SetupToken // token -> setup token
	mutex       sync.RWMutex
}

// NewMemoryStorage creates a new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		events:      make([]models.Event, 0),
		users:       make(map[string]*models.User),
		apiKeys:     make(map[string]*models.APIKey),
		setupTokens: make(map[string]*models.SetupToken),
	}

	// Add default user
	defaultUser, _ := models.NewUser("user-123")
	storage.SaveUser(defaultUser)

	// Add default API key for test user
	plainKey := "sk_testkey123456789012345678901234567890"
	keyHash, _ := bcrypt.GenerateFromPassword([]byte(plainKey), bcrypt.DefaultCost)
	now := time.Now()
	apiKey := &models.APIKey{
		UUID:         "test-api-key-uuid",
		UserID:       "user-123",
		KeyHash:      string(keyHash),
		EncryptedKey: "encrypted-test-key", // Not used in tests
		Description:  "Test API Key",
		CreatedAt:    now,
		LastUsedAt:   &now,
	}
	storage.CreateAPIKey(apiKey)

	return storage
}

// SaveEvents appends new events to the storage
func (m *MemoryStorage) SaveEvents(events []models.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.events = append(m.events, events...)
	return nil
}

// LoadEvents returns all stored events
func (m *MemoryStorage) LoadEvents() ([]models.Event, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external modification
	allEvents := make([]models.Event, len(m.events))
	copy(allEvents, m.events)

	return allEvents, nil
}

// SaveUser stores a user by id
func (m *MemoryStorage) SaveUser(user *models.User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.users[user.Id] = user
	return nil
}

// GetUserById retrieves a user by id
func (m *MemoryStorage) GetUserById(id string) (*models.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// CreateAPIKey stores a new API key
func (m *MemoryStorage) CreateAPIKey(apiKey *models.APIKey) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.apiKeys[apiKey.UUID] = apiKey
	return nil
}

// GetAPIKeyByHash retrieves an API key by its hash
func (m *MemoryStorage) GetAPIKeyByHash(hash string) (*models.APIKey, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, apiKey := range m.apiKeys {
		if apiKey.KeyHash == hash {
			return apiKey, nil
		}
	}
	return nil, errors.New("API key not found")
}

// GetAllAPIKeys retrieves all API keys
func (m *MemoryStorage) GetAllAPIKeys() ([]*models.APIKey, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	keys := make([]*models.APIKey, 0, len(m.apiKeys))
	for _, k := range m.apiKeys {
		keys = append(keys, k)
	}
	return keys, nil
}

// UpdateAPIKey updates an existing API key
func (m *MemoryStorage) UpdateAPIKey(apiKey *models.APIKey) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.apiKeys[apiKey.UUID] = apiKey
	return nil
}

// CreateSetupToken stores a new setup token
func (m *MemoryStorage) CreateSetupToken(token *models.SetupToken) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.setupTokens[token.Token] = token
	return nil
}

// GetSetupToken retrieves a setup token by its value
func (m *MemoryStorage) GetSetupToken(token string) (*models.SetupToken, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	setupToken, exists := m.setupTokens[token]
	if !exists {
		return nil, errors.New("setup token not found")
	}
	return setupToken, nil
}

// UpdateSetupToken updates an existing setup token
func (m *MemoryStorage) UpdateSetupToken(token *models.SetupToken) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.setupTokens[token.Token] = token
	return nil
}

// InvalidateUserSetupTokens marks all setup tokens for a user as used
func (m *MemoryStorage) InvalidateUserSetupTokens(userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, token := range m.setupTokens {
		if token.UserID == userID {
			token.Used = true
		}
	}
	return nil
}

// InvalidateUserAPIKeys removes all API keys for a user
func (m *MemoryStorage) InvalidateUserAPIKeys(userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for uuid, apiKey := range m.apiKeys {
		if apiKey.UserID == userID {
			delete(m.apiKeys, uuid)
		}
	}
	return nil
}
