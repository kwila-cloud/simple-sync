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
