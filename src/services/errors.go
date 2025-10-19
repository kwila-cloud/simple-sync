package services

import "simple-sync/src/errors"

// Re-export shared errors for convenience in services package
var (
	// Authentication errors
	ErrInvalidApiKeyFormat = errors.ErrInvalidApiKeyFormat
	ErrInvalidApiKey       = errors.ErrInvalidApiKey
	ErrInvalidSetupToken   = errors.ErrInvalidSetupToken
	ErrSetupTokenExpired   = errors.ErrSetupTokenExpired

	// Business logic errors
	ErrUserNotFound = errors.ErrUserNotFound
)
