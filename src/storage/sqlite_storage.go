package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"simple-sync/src/models"
)

// SQLiteStorage implements Storage using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates an instance
func NewSQLiteStorage() *SQLiteStorage {
	return &SQLiteStorage{}
}

// Initialize opens a connection to the SQLite database at path
func (s *SQLiteStorage) Initialize(path string) error {
	if path == "" {
		path = getDefaultDBPath()
	}

	// Determine and create parent directory unless using an in-memory DB
	if path != "" {
		isMemory := path == ":memory:" || (strings.HasPrefix(path, "file:") && strings.Contains(path, "memory"))
		if !isMemory {
			parent := filepath.Dir(path)
			if parent != "" && parent != "." {
				if err := os.MkdirAll(parent, 0o755); err != nil {
					return fmt.Errorf("failed to create db dir: %w", err)
				}
			}
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	// Set pragmas for safety/performance
	// Use WAL (Write-Ahead Logging) to improve concurrency and crash resilience:
	// - Allows readers to run while a writer is writing
	// - Better write throughput for many workloads
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return err
	}
	// Enforce foreign key constraints at the SQLite level to maintain relational integrity.
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		db.Close()
		return err
	}
	// Set synchronous to NORMAL to balance durability and performance:
	// - FULL is the most durable (safer on power loss) but slower
	// - NORMAL offers a good trade-off for many server environments
	if _, err := db.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		db.Close()
		return err
	}

	// Configure connection pool defaults
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	// Apply migrations (idempotent)
	if err := applyMigrations(db); err != nil {
		db.Close()
		return err
	}

	s.db = db
	return nil
}

// applyMigrations ensures required tables and indexes exist
func applyMigrations(db *sql.DB) error {
	stmts := []string{
		// users table
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			created_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);`,
		// events table
		`CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			item TEXT NOT NULL,
			type TEXT NOT NULL,
			data TEXT,
			created_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);`,
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

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	if s.db == nil {
		return nil
	}
	// Close DB and clear pointer
	err := s.db.Close()
	s.db = nil
	return err
}

// Minimal stubs to satisfy the Storage interface — to be implemented later
func (s *SQLiteStorage) SaveEvents(events []models.Event) error      { return ErrInvalidData }
func (s *SQLiteStorage) LoadEvents() ([]models.Event, error)         { return nil, ErrNotFound }
func (s *SQLiteStorage) SaveUser(user *models.User) error            { return ErrInvalidData }
func (s *SQLiteStorage) GetUserById(id string) (*models.User, error) { return nil, ErrNotFound }
func (s *SQLiteStorage) CreateApiKey(apiKey *models.ApiKey) error    { return ErrInvalidData }
func (s *SQLiteStorage) GetApiKeyByHash(hash string) (*models.ApiKey, error) {
	return nil, ErrApiKeyNotFound
}
func (s *SQLiteStorage) GetAllApiKeys() ([]*models.ApiKey, error)        { return nil, ErrNotFound }
func (s *SQLiteStorage) UpdateApiKey(apiKey *models.ApiKey) error        { return ErrInvalidData }
func (s *SQLiteStorage) InvalidateUserApiKeys(userID string) error       { return nil }
func (s *SQLiteStorage) CreateSetupToken(token *models.SetupToken) error { return ErrInvalidData }
func (s *SQLiteStorage) GetSetupToken(token string) (*models.SetupToken, error) {
	return nil, ErrSetupTokenNotFound
}
func (s *SQLiteStorage) UpdateSetupToken(token *models.SetupToken) error { return ErrInvalidData }
func (s *SQLiteStorage) InvalidateUserSetupTokens(userID string) error   { return nil }
func (s *SQLiteStorage) CreateAclRule(rule *models.AclRule) error        { return ErrInvalidData }
func (s *SQLiteStorage) GetAclRules() ([]models.AclRule, error)          { return nil, ErrNotFound }

func getDefaultDBPath() string {
	if p := os.Getenv("DB_PATH"); p != "" {
		return p
	}
	return "./data/simple-sync.db"
}
