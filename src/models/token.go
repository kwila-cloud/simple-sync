package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTToken represents an authentication token issued to a user
type JWTToken struct {
	TokenString string    `json:"token_string"`
	UserUUID    string    `json:"user_uuid"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsRevoked   bool      `json:"is_revoked"`
}

// Validate performs validation on the JWTToken struct
func (t *JWTToken) Validate() error {
	if t.TokenString == "" {
		return errors.New("token string is required")
	}

	if t.UserUUID == "" {
		return errors.New("user UUID is required")
	}

	if t.IssuedAt.IsZero() {
		return errors.New("issued at time is required")
	}

	if t.ExpiresAt.IsZero() {
		return errors.New("expires at time is required")
	}

	if t.ExpiresAt.Before(t.IssuedAt) {
		return errors.New("expires at must be after issued at")
	}

	return nil
}

// IsExpired checks if the token has expired
func (t *JWTToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// TokenClaims represents the payload embedded in JWT tokens
type TokenClaims struct {
	UserUUID  string `json:"user_uuid"`
	Username  string `json:"username"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	IsAdmin   bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Validate performs validation on the TokenClaims struct
func (c *TokenClaims) Validate() error {
	if c.UserUUID == "" {
		return errors.New("user UUID is required")
	}

	if c.Username == "" {
		return errors.New("username is required")
	}

	if c.IssuedAt == 0 {
		return errors.New("issued at is required")
	}

	if c.ExpiresAt == 0 {
		return errors.New("expires at is required")
	}

	if c.ExpiresAt <= c.IssuedAt {
		return errors.New("expires at must be after issued at")
	}

	return nil
}

// NewTokenClaims creates new token claims
func NewTokenClaims(userUUID, username string, isAdmin bool, issuedAt, expiresAt time.Time) *TokenClaims {
	return &TokenClaims{
		UserUUID:  userUUID,
		Username:  username,
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: expiresAt.Unix(),
		IsAdmin:   isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
}
