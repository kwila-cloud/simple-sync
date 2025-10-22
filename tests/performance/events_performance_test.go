package performance

import (
	"testing"

	"simple-sync/src/storage"
)

// BenchmarkSaveLoadEvents benchmarks saving and loading events with a large dataset
func BenchmarkSaveLoadEvents(b *testing.B) {
	// use helper to produce dataset
	events := GenerateEvents(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := storage.NewSQLiteStorage()
		if err := s.Initialize(":memory:"); err != nil {
			b.Fatalf("failed to init sqlite: %v", err)
		}
		// save
		if err := s.SaveEvents(events); err != nil {
			b.Fatalf("save events failed: %v", err)
		}
		// load
		if _, err := s.LoadEvents(); err != nil {
			b.Fatalf("load events failed: %v", err)
		}
		s.Close()
	}
}
