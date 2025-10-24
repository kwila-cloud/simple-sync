package performance

import (
	"sync"
	"testing"

	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

// TestSQLiteStorageConcurrent performs concurrent reads and writes against SQLiteStorage
func TestSQLiteStorageConcurrent(t *testing.T) {
	s := storage.NewSQLiteStorage()
	tmp := t.TempDir()
	dbPath := tmp + "/perf_test.db"
	if err := s.Initialize(dbPath); err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer s.Close()

	// Writers will each insert 500 events; readers will continuously read
	writers := 10
	eventsPerWriter := 500
	readers := 5

	var wg sync.WaitGroup
	errCh := make(chan error, writers+readers)

	// Start readers
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < writers; i++ {
				if _, err := s.LoadEvents(); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	// Start writers
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			events := make([]models.Event, 0, eventsPerWriter)
			for i := 0; i < eventsPerWriter; i++ {
				// Use NewEvent to ensure valid UUID/timestamp
				e := models.NewEvent("concurrent-user", "item", "action", "payload")
				events = append(events, *e)
			}
			if err := s.AddEvents(events); err != nil {
				errCh <- err
				return
			}
		}(w)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		assert.NoError(t, err)
	}

	// Verify total events written equals writers * eventsPerWriter
	events, err := s.LoadEvents()
	if err != nil {
		t.Fatalf("failed to load events: %v", err)
	}
	expected := writers * eventsPerWriter
	assert.Equal(t, expected, len(events), "expected total events to equal writers*eventsPerWriter")
}
