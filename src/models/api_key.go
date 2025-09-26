package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// APIKey represents a long-lived API key for user authentication
type APIKey struct {
	UUID        string     `json:"uuid" db:"uuid"`
	UserID      string     `json:"user_id" db:"user_id"`
	KeyHash     string     `json:"key_hash" db:"key_hash"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	Description string     `json:"description,omitempty" db:"description"`
}

// Validate performs validation on the APIKey struct
func (k *APIKey) Validate() error {
	if k.UUID == "" {
		return errors.New("UUID is required")
	}

	if _, err := uuid.Parse(k.UUID); err != nil {
		return errors.New("UUID must be valid format")
	}

	if k.UserID == "" {
		return errors.New("user ID is required")
	}

	if k.KeyHash == "" {
		return errors.New("key hash is required")
	}

	if k.CreatedAt.IsZero() {
		return errors.New("created at time is required")
	}

	return nil
}

// NewAPIKey creates a new API key instance
func NewAPIKey(userID, keyHash, description string) *APIKey {
	return &APIKey{
		UUID:        uuid.New().String(),
		UserID:      userID,
		KeyHash:     keyHash,
		CreatedAt:   time.Now(),
		Description: description,
	}
}

// UpdateLastUsed updates the last used timestamp
func (k *APIKey) UpdateLastUsed() {
	now := time.Now()
	k.LastUsedAt = &now
}
