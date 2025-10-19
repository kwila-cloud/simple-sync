package services

import "errors"

// Service-specific error types
var (
	// Authentication errors
	ErrInvalidApiKeyFormat = errors.New("invalid API key format")
	ErrInvalidApiKey       = errors.New("invalid API key")
	ErrInvalidSetupToken   = errors.New("invalid setup token")
	ErrSetupTokenExpired   = errors.New("setup token is expired or already used")

	// Business logic errors
	ErrUserNotFound = errors.New("user not found")
)
