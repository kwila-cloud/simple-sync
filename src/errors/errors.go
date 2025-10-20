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
	ErrInvalidTimestamp   = errors.New("invalid timestamp")
	ErrInvalidUuidFormat  = errors.New("UUID must be valid format")
	ErrUuidRequired       = errors.New("UUID is required")
	ErrUserRequired       = errors.New("user is required")
	ErrItemRequired       = errors.New("item is required")
	ErrActionRequired     = errors.New("action is required")
	ErrKeyHashRequired    = errors.New("key hash is required")
	ErrCreatedAtRequired  = errors.New("created at time is required")
	ErrTokenRequired      = errors.New("token is required")
	ErrTokenInvalidFormat = errors.New("token must be in format XXXX-XXXX")
	ErrExpiresAtRequired  = errors.New("expires at time is required")
	ErrIdRequired         = errors.New("id is required")

	// ACL validation errors
	ErrInvalidAclType            = errors.New("type must be either 'allow' or 'deny'")
	ErrAclUserEmpty              = errors.New("user pattern cannot be empty")
	ErrAclItemEmpty              = errors.New("item pattern cannot be empty")
	ErrAclActionEmpty            = errors.New("action pattern cannot be empty")
	ErrAclUserControlChars       = errors.New("user pattern contains invalid control characters")
	ErrAclItemControlChars       = errors.New("item pattern contains invalid control characters")
	ErrAclActionControlChars     = errors.New("action pattern contains invalid control characters")
	ErrAclUserInvalidWildcards   = errors.New("user pattern can have at most one wildcard at the end")
	ErrAclItemInvalidWildcards   = errors.New("item pattern can have at most one wildcard at the end")
	ErrAclActionInvalidWildcards = errors.New("action pattern can have at most one wildcard at the end")

	// Business logic errors
	ErrUserNotFound = errors.New("user not found")
)
