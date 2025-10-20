package models

import (
	"time"

	"github.com/google/uuid"
	apperrors "simple-sync/src/errors"
)

// ApiKey represents a long-lived API key for user authentication
type ApiKey struct {
	UUID        string     `json:"uuid" db:"uuid"`
	UserID      string     `json:"user_id" db:"user_id"`
	KeyHash     string     `json:"key_hash" db:"key_hash"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	Description string     `json:"description,omitempty" db:"description"`
}

// Validate performs validation on the ApiKey struct
func (k *ApiKey) Validate() error {
	if k.UUID == "" {
		return apperrors.ErrUuidRequired
	}

	_, err := uuid.Parse(k.UUID)
	if err != nil {
		return apperrors.ErrInvalidUuidFormat
	}

	if k.UserID == "" {
		return apperrors.ErrUserIdRequired
	}

	if k.KeyHash == "" {
		return apperrors.ErrKeyHashRequired
	}

	if k.CreatedAt.IsZero() {
		return apperrors.ErrCreatedAtRequired
	}

	return nil
}

// NewApiKey creates a new API key instance
func NewApiKey(userID, keyHash, description string) *ApiKey {
	keyUuid, _ := uuid.NewV7()
	unixTimeSeconds, _ := keyUuid.Time().UnixTime()

	return &ApiKey{
		UUID:        keyUuid.String(),
		UserID:      userID,
		KeyHash:     keyHash,
		CreatedAt:   time.Unix(unixTimeSeconds, 0),
		Description: description,
	}
}

// UpdateLastUsed updates the last used timestamp
func (k *ApiKey) UpdateLastUsed() {
	now := time.Now()
	k.LastUsedAt = &now
}
