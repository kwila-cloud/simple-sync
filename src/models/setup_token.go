package models

import (
	"errors"
	"regexp"
	"time"
)

// SetupToken represents a short-lived token for initial user authentication setup
type SetupToken struct {
	Token     string    `json:"token" db:"token"`
	UserID    string    `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
}

// Validate performs validation on the SetupToken struct
func (t *SetupToken) Validate() error {
	if t.Token == "" {
		return errors.New("token is required")
	}

	// Validate token format: XXXX-XXXX
	tokenRegex := regexp.MustCompile(`^[A-Z0-9]{4}-[A-Z0-9]{4}$`)
	if !tokenRegex.MatchString(t.Token) {
		return errors.New("token must be in format XXXX-XXXX")
	}

	if t.UserID == "" {
		return errors.New("user ID is required")
	}

	if t.ExpiresAt.IsZero() {
		return errors.New("expires at time is required")
	}

	return nil
}

// IsExpired checks if the token has expired
func (t *SetupToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsValid checks if the token is valid for use
func (t *SetupToken) IsValid() bool {
	return !t.Used && !t.IsExpired()
}

// MarkUsed marks the token as used
func (t *SetupToken) MarkUsed() {
	t.Used = true
}

// NewSetupToken creates a new setup token instance
func NewSetupToken(token, userID string, expiresAt time.Time) *SetupToken {
	return &SetupToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Used:      false,
	}
}
