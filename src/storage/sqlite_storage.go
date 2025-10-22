package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"simple-sync/src/models"

	_ "github.com/mattn/go-sqlite3"
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
	if err := ApplyMigrations(db); err != nil {
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
func (s *SQLiteStorage) SaveEvents(events []models.Event) error {
	if s.db == nil {
		return ErrInvalidData
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO event (uuid, timestamp, user, item, action, payload) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, e := range events {
		if _, err := stmt.Exec(e.UUID, int64(e.Timestamp), e.User, e.Item, e.Action, e.Payload); err != nil {
			tx.Rollback()
			// Map sqlite unique/constraint errors to ErrDuplicateKey
			if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint failed") {
				return ErrDuplicateKey
			}
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (s *SQLiteStorage) LoadEvents() ([]models.Event, error) {
	if s.db == nil {
		return nil, ErrNotFound
	}
	rows, err := s.db.Query(`SELECT uuid, timestamp, user, item, action, payload FROM event ORDER BY timestamp ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []models.Event
	for rows.Next() {
		var e models.Event
		var ts int64
		if err := rows.Scan(&e.UUID, &ts, &e.User, &e.Item, &e.Action, &e.Payload); err != nil {
			return nil, err
		}
		e.Timestamp = uint64(ts)
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
func (s *SQLiteStorage) AddUser(user *models.User) error {
	if s.db == nil {
		return ErrInvalidData
	}
	if user == nil {
		return ErrInvalidData
	}

	// Validate the user model
	if err := user.Validate(); err != nil {
		return err
	}

	// Insert user into the database
	_, err := s.db.Exec(`INSERT INTO user (id, created_at) VALUES (?, ?)`, user.Id, user.CreatedAt)
	if err != nil {
		// Map sqlite unique/constraint errors to ErrDuplicateKey
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint failed") {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}
func (s *SQLiteStorage) GetUserById(id string) (*models.User, error) {
	if s.db == nil {
		return nil, ErrNotFound
	}

	var uid string
	var createdAt time.Time
	row := s.db.QueryRow(`SELECT id, created_at FROM user WHERE id = ?`, id)
	if err := row.Scan(&uid, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user := &models.User{
		Id:        uid,
		CreatedAt: createdAt,
	}
	return user, nil
}
func (s *SQLiteStorage) CreateApiKey(apiKey *models.ApiKey) error { return ErrInvalidData }
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
func (s *SQLiteStorage) CreateAclRule(rule *models.AclRule) error {
	if s.db == nil || rule == nil {
		return ErrInvalidData
	}

	// Validate the rule
	if err := rule.Validate(); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`INSERT INTO acl_rule (user, item, action, type) VALUES (?, ?, ?, ?)`,
		rule.User, rule.Item, rule.Action, rule.Type)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint failed") {
			return ErrDuplicateKey
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
func (s *SQLiteStorage) GetAclRules() ([]models.AclRule, error) {
	if s.db == nil {
		return nil, ErrNotFound
	}

	rows, err := s.db.Query(`SELECT user, item, action, type FROM acl_rule ORDER BY user, item, action, type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.AclRule
	for rows.Next() {
		var r models.AclRule
		if err := rows.Scan(&r.User, &r.Item, &r.Action, &r.Type); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func getDefaultDBPath() string {
	if p := os.Getenv("DB_PATH"); p != "" {
		return p
	}
	return "./data/simple-sync.db"
}
