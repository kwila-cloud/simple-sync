package performance

import (
	"testing"
	"time"

	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

// TestSaveLoadEventsPerformance measures save+load time for 10k events
func TestSaveLoadEventsPerformance(t *testing.T) {
	events := GenerateEvents(10000)

	s := storage.NewSQLiteStorage()
	tmp := t.TempDir()
	dbPath := tmp + "/perf_save_load.db"
	if err := s.Initialize(dbPath); err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer s.Close()

	start := time.Now()
	if err := s.SaveEvents(events); err != nil {
		t.Fatalf("save events failed: %v", err)
	}
	if _, err := s.LoadEvents(); err != nil {
		t.Fatalf("load events failed: %v", err)
	}
	d := time.Since(start)
	// Target: complete save+load for 10k events in under 2s
	assert.Less(t, d, 2*time.Second, "Save+Load for 10k events should complete in under 2s")
}
