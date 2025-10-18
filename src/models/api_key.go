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

	uuid, err := uuid.Parse(k.UUID)
	if err != nil {
		return errors.New("UUID must be valid format")
	}

	timestamp, _ := uuid.Time().UnixTime()
	if timestamp != k.CreatedAt.Unix() {
		return errors.New("UUID must match timestamp")
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
	keyUuid, _ := uuid.NewV7()
	unixTimeSeconds, _ := keyUuid.Time().UnixTime()

	return &APIKey{
		UUID:        keyUuid.String(),
		UserID:      userID,
		KeyHash:     keyHash,
		CreatedAt:   time.Unix(unixTimeSeconds, 0),
		Description: description,
	}
}

// UpdateLastUsed updates the last used timestamp
func (k *APIKey) UpdateLastUsed() {
	now := time.Now()
	k.LastUsedAt = &now
}
