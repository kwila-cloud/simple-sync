package performance

import (
	"testing"

	"simple-sync/src/storage"
)

// Consolidated benchmarks for save+load of events (10k dataset)
// - BenchmarkSaveLoadEvents_Memory_10k measures an in-memory SQLite run (":memory:")
//   and includes initialization overhead. Useful for profiling pure in-memory DB performance.
// - BenchmarkSaveLoadEvents_Disk_10k measures a file-backed SQLite run and isolates the
//   save+load timing by stopping the timer during setup. Useful for realistic disk-backed measurements.

func BenchmarkSaveLoadEvents_Memory_10k(b *testing.B) {
	events := GenerateEvents(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := storage.NewSQLiteStorage()
		if err := s.Initialize(":memory:"); err != nil {
			b.Fatalf("failed to init sqlite: %v", err)
		}
		if err := s.SaveEvents(events); err != nil {
			b.Fatalf("save events failed: %v", err)
		}
		if _, err := s.LoadEvents(); err != nil {
			b.Fatalf("load events failed: %v", err)
		}
		s.Close()
	}
}

func BenchmarkSaveLoadEvents_Disk_10k(b *testing.B) {
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
