package storage

import (
	"testing"
)

func TestSQLiteMigrations(t *testing.T) {
	store := NewSQLiteStorage()
	if err := store.Initialize(":memory:"); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer store.Close()

	tables := []string{"users", "events", "api_keys", "setup_tokens", "acl_rules"}
	for _, table := range tables {
		var name string
		err := store.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Fatalf("expected table %s exists: %v", table, err)
		}
		if name != table {
			t.Fatalf("expected table %s, got %s", table, name)
		}
	}
}
