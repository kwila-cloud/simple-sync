package storage

import (
	"sync"

	"simple-sync/src/models"
)

// MemoryStorage implements in-memory storage for events
type MemoryStorage struct {
	events []models.Event
	mutex  sync.RWMutex
}

// NewMemoryStorage creates a new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		events: make([]models.Event, 0),
	}
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
	events := make([]models.Event, len(m.events))
	copy(events, m.events)
	return events, nil
}