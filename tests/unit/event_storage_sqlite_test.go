package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestAddAndLoadEvents(t *testing.T) {
	s := storage.NewSQLiteStorage()
	if err := s.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize in-memory sqlite: %v", err)
	}
	defer s.Close()

	e1 := models.NewEvent("user1", "item1", "act1", "payload1")
	// small delay to ensure timestamp ordering
	time.Sleep(1 * time.Millisecond)
	e2 := models.NewEvent("user2", "item2", "act2", "payload2")

	if err := s.AddEvents([]models.Event{*e1, *e2}); err != nil {
		t.Fatalf("AddEvents failed: %v", err)
	}

	events, err := s.LoadEvents()
	if err != nil {
		t.Fatalf("LoadEvents failed: %v", err)
	}

	assert.Equal(t, 2, len(events))
	assert.Equal(t, e1.UUID, events[0].UUID)
	assert.Equal(t, e2.UUID, events[1].UUID)
}

func TestAddEventsDuplicateUUID(t *testing.T) {
	s := storage.NewSQLiteStorage()
	if err := s.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize in-memory sqlite: %v", err)
	}
	defer s.Close()

	e := models.NewEvent("user1", "item1", "act1", "payload1")
	if err := s.AddEvents([]models.Event{*e}); err != nil {
		t.Fatalf("first AddEvents failed: %v", err)
	}
	// attempt to insert duplicate
	if err := s.AddEvents([]models.Event{*e}); err == nil {
		t.Fatalf("expected duplicate AddEvents to fail")
	} else {
		assert.Equal(t, storage.ErrDuplicateKey, err)
	}
}

func TestLoadEventsEmpty(t *testing.T) {
	s := storage.NewSQLiteStorage()
	if err := s.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize in-memory sqlite: %v", err)
	}
	defer s.Close()

	events, err := s.LoadEvents()
	if err != nil {
		t.Fatalf("LoadEvents failed: %v", err)
	}
	assert.Equal(t, 0, len(events))
}
