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

	s.db = db
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

func getDir(path string) string {
	// Empty path: nothing to create here (caller will use default path)
	if path == "" {
		return ""
	}

	// Handle in-memory SQLite URIs (e.g. ":memory:", "file::memory:?cache=shared", or any
	// file: URI that contains "memory"). For these we don't create any directory.
	if path == ":memory:" || strings.HasPrefix(path, "file:") && strings.Contains(path, "memory") {
		return ""
	}

	// Return the parent directory for the given path (e.g. "./data/file.db" -> "./data").
	// filepath.Dir returns "." for paths without a directory component; callers can
	// safely call MkdirAll on "." (no-op) but we return the raw Dir value for clarity.
	return filepath.Dir(path)
}
