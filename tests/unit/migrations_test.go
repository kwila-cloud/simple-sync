package unit

import (
	"database/sql"
	"simple-sync/src/storage"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestApplyMigrationsCreatesTablesAndSetsVersion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	defer db.Close()

	if err := storage.ApplyMigrations(db); err != nil {
		t.Fatalf("ApplyMigrations failed: %v", err)
	}

	tables := []string{"user", "event", "api_key", "setup_token", "acl_rule"}
	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Fatalf("expected table %s exists: %v", table, err)
		}
		if name != table {
			t.Fatalf("expected table %s, got %s", table, name)
		}
	}

	var v int
	if err := db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		t.Fatalf("failed to read user_version: %v", err)
	}
	if v != storage.DesiredSchemaVersion {
		t.Fatalf("expected user_version %d, got %d", storage.DesiredSchemaVersion, v)
	}

	// Idempotent: second run should succeed and leave version unchanged
	if err := storage.ApplyMigrations(db); err != nil {
		t.Fatalf("second ApplyMigrations failed: %v", err)
	}
	if err := db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		t.Fatalf("failed to read user_version after second run: %v", err)
	}
	if v != storage.DesiredSchemaVersion {
		t.Fatalf("expected user_version %d after second run, got %d", storage.DesiredSchemaVersion, v)
	}
}
