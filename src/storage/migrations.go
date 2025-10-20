package storage

import "database/sql"

// ApplyMigrations ensures required tables and indexes exist (exported for testing)
func ApplyMigrations(db *sql.DB) error {
	stmts := []string{
		// users table
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			created_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);`,
		// events table - match models.Event
		`CREATE TABLE IF NOT EXISTS events (
			uuid TEXT PRIMARY KEY,
			timestamp INTEGER NOT NULL,
			user TEXT NOT NULL,
			item TEXT NOT NULL,
			action TEXT NOT NULL,
			payload TEXT,
			created_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);`,
		`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);`,
		`CREATE INDEX IF NOT EXISTS idx_events_item ON events(item);`,
		// api_keys table
		`CREATE TABLE IF NOT EXISTS api_keys (
			uuid TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			key_hash TEXT NOT NULL,
			description TEXT,
			created_at DATETIME NOT NULL DEFAULT (datetime('now')),
			last_used_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);`,
		// setup_tokens table
		`CREATE TABLE IF NOT EXISTS setup_tokens (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at DATETIME,
			used_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		// acl_rules table
		`CREATE TABLE IF NOT EXISTS acl_rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			subject TEXT NOT NULL,
			object TEXT NOT NULL,
			action TEXT NOT NULL,
			meta TEXT,
			created_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);`,
		`CREATE INDEX IF NOT EXISTS idx_acl_subject_object ON acl_rules(subject, object);`,
	}

	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
