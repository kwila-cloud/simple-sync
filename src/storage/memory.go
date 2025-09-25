package storage

import (
	"errors"
	"sync"

	"simple-sync/src/models"
)

// MemoryStorage implements in-memory storage
type MemoryStorage struct {
	events      []models.Event
	users       map[string]*models.User       // username -> user
	apiKeys     map[string]*models.APIKey     // hash -> api key
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
	defaultUser, _ := models.NewUserWithPassword("user-123", "testuser", "testpass123", false)
	storage.SaveUser(defaultUser)

	return storage
}

// SaveEvents appends new events to the storage
func (m *MemoryStorage) SaveEvents(events []models.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.events = append(m.events, events...)
	return nil
}

// LoadEvents returns all stored events, optionally filtered by fromTimestamp
func (m *MemoryStorage) LoadEvents(fromTimestamp *uint64) ([]models.Event, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	allEvents := make([]models.Event, len(m.events))
	copy(allEvents, m.events)

	// If no timestamp filter is provided, return all events
	if fromTimestamp == nil {
		return allEvents, nil
	}

	// Filter events by timestamp
	filteredEvents := make([]models.Event, 0)
	for _, event := range allEvents {
		if event.Timestamp >= *fromTimestamp {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}

// SaveUser stores a user by username
func (m *MemoryStorage) SaveUser(user *models.User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.users[user.Username] = user
	return nil
}

// GetUserByUsername retrieves a user by username
func (m *MemoryStorage) GetUserByUsername(username string) (*models.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	user, exists := m.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByUUID retrieves a user by UUID
func (m *MemoryStorage) GetUserByUUID(uuid string) (*models.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, user := range m.users {
		if user.UUID == uuid {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// CreateAPIKey stores a new API key
func (m *MemoryStorage) CreateAPIKey(apiKey *models.APIKey) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.apiKeys[apiKey.KeyHash] = apiKey
	return nil
}

// GetAPIKeyByHash retrieves an API key by its hash
func (m *MemoryStorage) GetAPIKeyByHash(hash string) (*models.APIKey, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	apiKey, exists := m.apiKeys[hash]
	if !exists {
		return nil, errors.New("API key not found")
	}
	return apiKey, nil
}

// UpdateAPIKey updates an existing API key
func (m *MemoryStorage) UpdateAPIKey(apiKey *models.APIKey) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.apiKeys[apiKey.KeyHash] = apiKey
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
