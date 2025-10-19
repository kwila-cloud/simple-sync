package handlers

import "simple-sync/src/errors"

// Re-export shared errors for convenience in handlers package
var (
	// Authentication errors
	ErrInvalidApiKeyFormat = errors.ErrInvalidApiKeyFormat
	ErrInvalidApiKey       = errors.ErrInvalidApiKey
	ErrInvalidSetupToken   = errors.ErrInvalidSetupToken
	ErrSetupTokenExpired   = errors.ErrSetupTokenExpired

	// Validation errors
	ErrInvalidTimestamp = errors.ErrInvalidTimestamp

	// ACL validation errors
	ErrInvalidAclType             = errors.ErrInvalidAclType
	ErrAclUserEmpty               = errors.ErrAclUserEmpty
	ErrAclItemEmpty               = errors.ErrAclItemEmpty
	ErrAclActionEmpty             = errors.ErrAclActionEmpty
	ErrAclUserMultipleWildcards   = errors.ErrAclUserMultipleWildcards
	ErrAclItemMultipleWildcards   = errors.ErrAclItemMultipleWildcards
	ErrAclActionMultipleWildcards = errors.ErrAclActionMultipleWildcards
	ErrAclUserControlChars        = errors.ErrAclUserControlChars
	ErrAclItemControlChars        = errors.ErrAclItemControlChars
	ErrAclActionControlChars      = errors.ErrAclActionControlChars
)
