package storage

import (
	"database/sql"
	"fmt"
	"log"
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
		var err error
		path, err = getDbPath()
		if err != nil {
			return fmt.Errorf("failed to get database path")
		}
	}
	log.Printf("initializing with path: %v", path)

	// Determine and create parent directory unless using an in-memory DB
	if path != "" {
		isMemory := path == ":memory:" || (strings.HasPrefix(path, "file:") && strings.Contains(path, "memory"))
		if !isMemory {
			parent := filepath.Dir(path)
			if parent != "" && parent != "." {
				if err := os.MkdirAll(parent, 0o755); err != nil {
					return fmt.Errorf("failed to create database directory: %w", err)
				}
			}
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open SQLite file: %v, error: %v", path, err)
	}
	log.Printf("successfully opened SQLite file")

	// Set pragmas for safety/performance
	// Use WAL (Write-Ahead Logging) to improve concurrency and crash resilience:
	// - Allows readers to run while a writer is writing
	// - Better write throughput for many workloads
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		defer db.Close()
		return fmt.Errorf("failed to set journal mode: %v", err)
	}
	log.Println("A")
	// Enforce foreign key constraints at the SQLite level to maintain relational integrity.
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		defer db.Close()
		return fmt.Errorf("failed to enable foreign keys: %v", err)
	}

	log.Println("B")
	// Configure connection pool defaults
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("C")
	// Verify connection
	if err := db.Ping(); err != nil {
		defer db.Close()
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Apply migrations (idempotent)
	if err := ApplyMigrations(db); err != nil {
		defer db.Close()
		return fmt.Errorf("failed to apply migrations: %v", err)
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
	return fmt.Errorf("failed to close db: %v", err)
}

func (s *SQLiteStorage) AddEvents(events []models.Event) error {
	if s.db == nil {
		return ErrInvalidData
	}
	// Validate each event before attempting DB operations
	for i := range events {
		if err := events[i].Validate(); err != nil {
			return err
		}
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
func (s *SQLiteStorage) AddApiKey(apiKey *models.ApiKey) error {
	if s.db == nil || apiKey == nil {
		return ErrInvalidData
	}
	if err := apiKey.Validate(); err != nil {
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

	_, err = tx.Exec(`INSERT INTO api_key (uuid, user, key_hash, created_at, last_used_at, description) VALUES (?, ?, ?, ?, ?, ?)`,
		apiKey.UUID, apiKey.User, apiKey.KeyHash, apiKey.CreatedAt, apiKey.LastUsedAt, apiKey.Description)
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
func (s *SQLiteStorage) GetApiKeyByHash(hash string) (*models.ApiKey, error) {
	if s.db == nil {
		return nil, ErrApiKeyNotFound
	}
	row := s.db.QueryRow(`SELECT uuid, user, key_hash, created_at, last_used_at, description FROM api_key WHERE key_hash = ?`, hash)
	var k models.ApiKey
	var lastUsed sql.NullTime
	if err := row.Scan(&k.UUID, &k.User, &k.KeyHash, &k.CreatedAt, &lastUsed, &k.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrApiKeyNotFound
		}
		return nil, err
	}
	if lastUsed.Valid {
		k.LastUsedAt = &lastUsed.Time
	}
	return &k, nil
}
func (s *SQLiteStorage) GetAllApiKeys() ([]*models.ApiKey, error) {
	if s.db == nil {
		return nil, ErrNotFound
	}
	rows, err := s.db.Query(`SELECT uuid, user, key_hash, created_at, last_used_at, description FROM api_key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*models.ApiKey
	for rows.Next() {
		var k models.ApiKey
		var lastUsed sql.NullTime
		if err := rows.Scan(&k.UUID, &k.User, &k.KeyHash, &k.CreatedAt, &lastUsed, &k.Description); err != nil {
			return nil, err
		}
		if lastUsed.Valid {
			k.LastUsedAt = &lastUsed.Time
		}
		keys = append(keys, &k)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}
func (s *SQLiteStorage) UpdateApiKey(apiKey *models.ApiKey) error {
	if s.db == nil || apiKey == nil {
		return ErrInvalidData
	}
	// Only update fields that can change: last_used_at, description
	_, err := s.db.Exec(`UPDATE api_key SET last_used_at = ?, description = ? WHERE uuid = ?`, apiKey.LastUsedAt, apiKey.Description, apiKey.UUID)
	if err != nil {
		return err
	}
	return nil
}
func (s *SQLiteStorage) InvalidateUserApiKeys(userID string) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`DELETE FROM api_key WHERE user = ?`, userID)
	return err
}
func (s *SQLiteStorage) AddSetupToken(token *models.SetupToken) error {
	if s.db == nil || token == nil {
		return ErrInvalidData
	}
	if err := token.Validate(); err != nil {
		return err
	}
	_, err := s.db.Exec(`INSERT INTO setup_token (token, user, expires_at, used_at) VALUES (?, ?, ?, ?)`, token.Token, token.User, token.ExpiresAt, token.UsedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint failed") {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}
func (s *SQLiteStorage) GetSetupToken(token string) (*models.SetupToken, error) {
	if s.db == nil {
		return nil, ErrSetupTokenNotFound
	}
	row := s.db.QueryRow(`SELECT token, user, expires_at, used_at FROM setup_token WHERE token = ?`, token)
	var st models.SetupToken
	var used sql.NullTime
	if err := row.Scan(&st.Token, &st.User, &st.ExpiresAt, &used); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSetupTokenNotFound
		}
		return nil, err
	}
	if used.Valid {
		st.UsedAt = used.Time
	}
	return &st, nil
}
func (s *SQLiteStorage) UpdateSetupToken(token *models.SetupToken) error {
	if s.db == nil || token == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`UPDATE setup_token SET user = ?, expires_at = ?, used_at = ? WHERE token = ?`, token.User, token.ExpiresAt, token.UsedAt, token.Token)
	if err != nil {
		return err
	}
	return nil
}
func (s *SQLiteStorage) InvalidateUserSetupTokens(userID string) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`UPDATE setup_token SET used_at = ? WHERE user = ?`, time.Now(), userID)
	return err
}
func (s *SQLiteStorage) AddAclRule(rule *models.AclRule) error {
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

	rows, err := s.db.Query(`SELECT user, item, action, type FROM acl_rule`)
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

// Get the DB path from the DB_PATH environment variable, if it exists.
// Otherwise uses ./data/simple-sync.db
// Paths are returned as absolute paths for deterministic results.
func getDbPath() (string, error) {
	if p := os.Getenv("DB_PATH"); p != "" {
		abs, err := filepath.Abs(p)
		if err == nil {
			return abs, nil
		}
		return "", err
	}
	abs, err := filepath.Abs("./data/simple-sync.db")
	if err == nil {
		return abs, nil
	}
	return "", err
}
