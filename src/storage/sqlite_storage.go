package storage

import (
	"database/sql"
	"fmt"
	"os"

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
	// Ensure directory exists
	dir := getDir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create db dir: %w", err)
		}
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
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
	return s.db.Close()
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
	if path == ":memory:" || path == "file::memory:?cache=shared" {
		return ""
	}
	return os.DirFS(path).Name()
}
