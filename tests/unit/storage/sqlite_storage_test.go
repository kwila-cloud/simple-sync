package storage_test

import (
	"testing"

	"simple-sync/src/storage"
)

func TestInitializeClose(t *testing.T) {
	s := storage.NewSQLiteStorage()
	// use in-memory SQLite for tests
	if err := s.Initialize("file::memory:?cache=shared"); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}
