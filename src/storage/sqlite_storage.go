package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// Event operations
func (s *SQLiteStorage) SaveEvents(events []models.Event) error {
	if s.db == nil {
		return ErrInvalidData
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
	stmt, err := tx.Prepare(`INSERT INTO events(item, type, data, created_at) VALUES(?, ?, ?, datetime('now'))`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, ev := range events {
		if _, err := stmt.Exec(ev.Item, ev.Action, ev.Payload); err != nil {
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
	rows, err := s.db.Query(`SELECT id, item, type, data, created_at FROM events ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []models.Event
	for rows.Next() {
		var id int
		var item, typ, data, createdAt string
		if err := rows.Scan(&id, &item, &typ, &data, &createdAt); err != nil {
			return nil, err
		}
		// Minimal mapping back to models.Event — UUID/timestamp mapping not preserved yet
		events = append(events, models.Event{UUID: "", User: "", Item: item, Action: typ, Payload: data})
	}
	return events, nil
}

// User operations
func (s *SQLiteStorage) SaveUser(user *models.User) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`INSERT OR REPLACE INTO users(id, created_at) VALUES(?, ?)`, user.Id, user.CreatedAt)
	return err
}

func (s *SQLiteStorage) GetUserById(id string) (*models.User, error) {
	if s.db == nil {
		return nil, ErrNotFound
	}
	row := s.db.QueryRow(`SELECT id, created_at FROM users WHERE id = ?`, id)
	var uid string
	var createdAt string
	if err := row.Scan(&uid, &createdAt); err != nil {
		return nil, ErrNotFound
	}
	// Parse createdAt
	return &models.User{Id: uid}, nil
}

// API Key operations
func (s *SQLiteStorage) CreateApiKey(apiKey *models.ApiKey) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`INSERT INTO api_keys(uuid, user_id, key_hash, description, created_at, last_used_at) VALUES(?, ?, ?, ?, ?, ?)`,
		apiKey.UUID, apiKey.User, apiKey.KeyHash, apiKey.Description, apiKey.CreatedAt, apiKey.LastUsedAt)
	return err
}

func (s *SQLiteStorage) GetApiKeyByHash(hash string) (*models.ApiKey, error) {
	if s.db == nil {
		return nil, ErrApiKeyNotFound
	}
	row := s.db.QueryRow(`SELECT uuid, user_id, key_hash, description, created_at, last_used_at FROM api_keys WHERE key_hash = ?`, hash)
	var uuid, userID, keyHash, description string
	var createdAt string
	var lastUsedAt sql.NullString
	if err := row.Scan(&uuid, &userID, &keyHash, &description, &createdAt, &lastUsedAt); err != nil {
		return nil, ErrApiKeyNotFound
	}
	// Minimal parsing of createdAt/lastUsedAt
	return &models.ApiKey{UUID: uuid, User: userID, KeyHash: keyHash, Description: description}, nil
}

func (s *SQLiteStorage) GetAllApiKeys() ([]*models.ApiKey, error) {
	if s.db == nil {
		return nil, ErrNotFound
	}
	rows, err := s.db.Query(`SELECT uuid, user_id, key_hash, description, created_at, last_used_at FROM api_keys`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var keys []*models.ApiKey
	for rows.Next() {
		var uuid, userID, keyHash, description string
		var createdAt string
		var lastUsedAt sql.NullString
		if err := rows.Scan(&uuid, &userID, &keyHash, &description, &createdAt, &lastUsedAt); err != nil {
			return nil, err
		}
		keys = append(keys, &models.ApiKey{UUID: uuid, User: userID, KeyHash: keyHash, Description: description})
	}
	return keys, nil
}

func (s *SQLiteStorage) UpdateApiKey(apiKey *models.ApiKey) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`UPDATE api_keys SET description = ?, last_used_at = ? WHERE uuid = ?`, apiKey.Description, apiKey.LastUsedAt, apiKey.UUID)
	return err
}

func (s *SQLiteStorage) InvalidateUserApiKeys(userID string) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`DELETE FROM api_keys WHERE user_id = ?`, userID)
	return err
}

// Setup Token operations
func (s *SQLiteStorage) CreateSetupToken(token *models.SetupToken) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`INSERT INTO setup_tokens(token, user_id, expires_at, used_at, created_at) VALUES(?, ?, ?, ?, ?)`,
		token.Token, token.User, token.ExpiresAt, token.UsedAt, time.Now())
	return err
}

func (s *SQLiteStorage) GetSetupToken(token string) (*models.SetupToken, error) {
	if s.db == nil {
		return nil, ErrSetupTokenNotFound
	}
	row := s.db.QueryRow(`SELECT token, user_id, expires_at, used_at, created_at FROM setup_tokens WHERE token = ?`, token)
	var tk, userID string
	var expiresAt sql.NullString
	var usedAt sql.NullString
	var createdAt string
	if err := row.Scan(&tk, &userID, &expiresAt, &usedAt, &createdAt); err != nil {
		return nil, ErrSetupTokenNotFound
	}
	return &models.SetupToken{Token: tk, User: userID}, nil
}

func (s *SQLiteStorage) UpdateSetupToken(token *models.SetupToken) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`UPDATE setup_tokens SET expires_at = ?, used_at = ? WHERE token = ?`, token.ExpiresAt, token.UsedAt, token.Token)
	return err
}

func (s *SQLiteStorage) InvalidateUserSetupTokens(userID string) error {
	if s.db == nil {
		return ErrInvalidData
	}
	_, err := s.db.Exec(`UPDATE setup_tokens SET used_at = datetime('now') WHERE user_id = ?`, userID)
	return err
}

// ACL operations
func (s *SQLiteStorage) CreateAclRule(rule *models.AclRule) error {
	if s.db == nil {
		return ErrInvalidData
	}
	ruleJson, _ := json.Marshal(rule)
	_, err := s.db.Exec(`INSERT INTO acl_rules(subject, object, action, meta, created_at) VALUES(?, ?, ?, ?, datetime('now'))`, rule.User, rule.Item, rule.Action, string(ruleJson))
	return err
}

func (s *SQLiteStorage) GetAclRules() ([]models.AclRule, error) {
	if s.db == nil {
		return nil, ErrNotFound
	}
	rows, err := s.db.Query(`SELECT meta FROM acl_rules ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []models.AclRule
	for rows.Next() {
		var meta string
		if err := rows.Scan(&meta); err != nil {
			return nil, err
		}
		var r models.AclRule
		if err := json.Unmarshal([]byte(meta), &r); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func getDefaultDBPath() string {
	if p := os.Getenv("DB_PATH"); p != "" {
		return p
	}
	return "./data/simple-sync.db"
}
