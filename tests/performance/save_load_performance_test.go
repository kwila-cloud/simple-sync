package performance

import (
	"testing"
	"time"

	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

// TestAddEventsPerformance measures save time for many events
func TestAddEventsPerformance(t *testing.T) {
	events := GenerateEvents(1_000_000)

	s := storage.NewSQLiteStorage()
	tmp := t.TempDir()
	dbPath := tmp + "/perf_save.db"
	if err := s.Initialize(dbPath); err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer s.Close()

	start := time.Now()
	if err := s.AddEvents(events); err != nil {
		t.Fatalf("save events failed: %v", err)
	}
	d := time.Since(start)
	assert.Less(t, d, 8*time.Second, "Adding 1 million events should complete in under 8s")
}

// TestLoadEventsPerformance measures load time for many events
func TestLoadEventsPerformance(t *testing.T) {
	events := GenerateEvents(1_000_000)

	s := storage.NewSQLiteStorage()
	tmp := t.TempDir()
	dbPath := tmp + "/perf_load.db"
	if err := s.Initialize(dbPath); err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer s.Close()

	// Pre-populate the DB (not measured)
	if err := s.AddEvents(events); err != nil {
		t.Fatalf("pre-populate save events failed: %v", err)
	}

	start := time.Now()
	if _, err := s.LoadEvents(); err != nil {
		t.Fatalf("load events failed: %v", err)
	}
	d := time.Since(start)
	assert.Less(t, d, 4*time.Second, "Loading 1 million events should complete in under 4s")
}
