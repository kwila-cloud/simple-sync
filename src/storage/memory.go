package storage

import (
	"errors"
	"sync"

	"simple-sync/src/models"
)

// MemoryStorage implements in-memory storage
type MemoryStorage struct {
	events []models.Event
	users  map[string]*models.User // username -> user
	mutex  sync.RWMutex
}

// NewMemoryStorage creates a new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		events: make([]models.Event, 0),
		users:  make(map[string]*models.User),
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
