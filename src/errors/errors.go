package errors

import "errors"

// Shared error definitions used across the application
var (
	// Authentication errors
	ErrInvalidApiKeyFormat = errors.New("invalid API key format")
	ErrInvalidApiKey       = errors.New("invalid API key")
	ErrInvalidSetupToken   = errors.New("invalid setup token")
	ErrSetupTokenExpired   = errors.New("setup token is expired or already used")

	// Validation errors
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// ACL validation errors
	ErrInvalidAclType             = errors.New("type must be either 'allow' or 'deny'")
	ErrAclUserEmpty               = errors.New("user pattern cannot be empty")
	ErrAclItemEmpty               = errors.New("item pattern cannot be empty")
	ErrAclActionEmpty             = errors.New("action pattern cannot be empty")
	ErrAclUserMultipleWildcards   = errors.New("user pattern can have at most one wildcard at the end")
	ErrAclItemMultipleWildcards   = errors.New("item pattern can have at most one wildcard at the end")
	ErrAclActionMultipleWildcards = errors.New("action pattern can have at most one wildcard at the end")
	ErrAclUserControlChars        = errors.New("user pattern contains invalid control characters")
	ErrAclItemControlChars        = errors.New("item pattern contains invalid control characters")
	ErrAclActionControlChars      = errors.New("action pattern contains invalid control characters")

	// Business logic errors
	ErrUserNotFound = errors.New("user not found")
)
