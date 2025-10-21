package storage

import (
	"database/sql"
	"fmt"
)

// DesiredSchemaVersion is the latest schema version the app expects.
const DesiredSchemaVersion = 1

// migrations holds per-version migration functions that bring the DB to that version.
var migrations = map[int]func(tx *sql.Tx) error{
	1: func(tx *sql.Tx) error {
		stmts := []string{
			// users table
			`CREATE TABLE IF NOT EXISTS users (
				id TEXT PRIMARY KEY,
				created_at DATETIME NOT NULL
			);`,
			// events table - for models.Event
			`CREATE TABLE IF NOT EXISTS events (
				uuid TEXT PRIMARY KEY,
				timestamp INTEGER NOT NULL,
				user TEXT NOT NULL,
				item TEXT NOT NULL,
				action TEXT NOT NULL,
				payload TEXT
			);`,
			`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);`,
			`CREATE INDEX IF NOT EXISTS idx_events_item ON events(item);`,
			// api_keys table - for models.ApiKey
			`CREATE TABLE IF NOT EXISTS api_keys (
				uuid TEXT PRIMARY KEY,
				user TEXT NOT NULL,
				key_hash TEXT NOT NULL,
				created_at DATETIME NOT NULL,
				last_used_at DATETIME,
				description TEXT,
				FOREIGN KEY(user) REFERENCES users(id) ON DELETE CASCADE
			);`,
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);`,
			// setup_tokens table - for models.SetupToken
			`CREATE TABLE IF NOT EXISTS setup_tokens (
				token TEXT PRIMARY KEY,
				user TEXT NOT NULL,
				expires_at DATETIME,
				used_at DATETIME,
				FOREIGN KEY(user) REFERENCES users(id) ON DELETE CASCADE
			);`,
			// acl_rules table - for models.AclRule
			`CREATE TABLE IF NOT EXISTS acl_rules (
				user TEXT NOT NULL,
				item TEXT NOT NULL,
				action TEXT NOT NULL,
				type TEXT NOT NULL
				PRIMARY KEY (user, item, action, type)
			);`,
		}

		for _, s := range stmts {
			if _, err := tx.Exec(s); err != nil {
				return err
			}
		}
		return nil
	},
}

func getUserVersion(db *sql.DB) (int, error) {
	row := db.QueryRow("PRAGMA user_version")
	var v int
	if err := row.Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func setUserVersion(tx *sql.Tx, v int) error {
	_, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d;", v))
	return err
}

// ApplyMigrations migrates the DB from its current user_version up to DesiredSchemaVersion.
func ApplyMigrations(db *sql.DB) error {
	cur, err := getUserVersion(db)
	if err != nil {
		return err
	}
	if cur > DesiredSchemaVersion {
		return fmt.Errorf("database version %d is newer than application supports (%d)", cur, DesiredSchemaVersion)
	}

	for v := cur + 1; v <= DesiredSchemaVersion; v++ {
		mig, ok := migrations[v]
		if !ok {
			return fmt.Errorf("missing migration for version %d", v)
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if err := mig(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d failed: %w", v, err)
		}

		if err := setUserVersion(tx, v); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
