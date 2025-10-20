package models

import (
	"strings"

	apperrors "simple-sync/src/errors"
)

// AclRule represents an access control rule
type AclRule struct {
	User   string `json:"user" db:"user"`
	Item   string `json:"item" db:"item"`
	Action string `json:"action" db:"action"`
	Type   string `json:"type" db:"type"`
}

// Validate performs comprehensive validation on the AclRule struct
func (r *AclRule) Validate() error {
	if err := r.validatePattern(r.User, "user"); err != nil {
		return err
	}
	if err := r.validatePattern(r.Item, "item"); err != nil {
		return err
	}
	if err := r.validatePattern(r.Action, "action"); err != nil {
		return err
	}
	if r.Type != "allow" && r.Type != "deny" {
		return apperrors.ErrInvalidAclType
	}
	return nil
}

// validatePattern validates a pattern with specific error messages based on the field name
func (r *AclRule) validatePattern(pattern, fieldType string) error {
	if pattern == "" {
		switch fieldType {
		case "user":
			return apperrors.ErrAclUserEmpty
		case "item":
			return apperrors.ErrAclItemEmpty
		case "action":
			return apperrors.ErrAclActionEmpty
		}
	}
	if pattern == "*" {
		return nil
	}
	if r.containsControlChars(pattern) {
		switch fieldType {
		case "user":
			return apperrors.ErrAclUserControlChars
		case "item":
			return apperrors.ErrAclItemControlChars
		case "action":
			return apperrors.ErrAclActionControlChars
		}
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		if strings.Contains(prefix, "*") {
			switch fieldType {
			case "user":
				return apperrors.ErrAclUserInvalidWildcards
			case "item":
				return apperrors.ErrAclItemInvalidWildcards
			case "action":
				return apperrors.ErrAclActionInvalidWildcards
			}
		}
	} else if strings.Contains(pattern, "*") {
		switch fieldType {
		case "user":
			return apperrors.ErrAclUserInvalidWildcards
		case "item":
			return apperrors.ErrAclItemInvalidWildcards
		case "action":
			return apperrors.ErrAclActionInvalidWildcards
		}
	}
	return nil
}

// containsControlChars checks if string contains control characters
func (r *AclRule) containsControlChars(s string) bool {
	for _, r := range s {
		if r < 32 || r == 127 {
			return true
		}
	}
	return false
}
