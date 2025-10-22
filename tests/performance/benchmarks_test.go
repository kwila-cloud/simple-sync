package performance

import (
	"testing"

	"simple-sync/src/storage"
)

// BenchmarkSaveLoadEvents10k benchmarks save+load for 10k events.
func BenchmarkSaveLoadEvents10k(b *testing.B) {
	events := GenerateEvents(10000)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := storage.NewSQLiteStorage()
		tmp := b.TempDir()
		dbPath := tmp + "/bench_10k.db"
		if err := s.Initialize(dbPath); err != nil {
			b.Fatalf("failed to init sqlite: %v", err)
		}
		b.StartTimer()
		if err := s.SaveEvents(events); err != nil {
			b.Fatalf("save events failed: %v", err)
		}
		if _, err := s.LoadEvents(); err != nil {
			b.Fatalf("load events failed: %v", err)
		}
		b.StopTimer()
		s.Close()
	}
}
