package models

import (
	"errors"
	"strings"
)

// ACLEvent represents an access control rule event
type ACLEvent struct {
	User   string `json:"user" binding:"required"`
	Item   string `json:"item" binding:"required"`
	Action string `json:"action" binding:"required"`
}

// Validate checks if the ACL event has valid data
func (a *ACLEvent) Validate() error {
	if strings.TrimSpace(a.User) == "" {
		return errors.New("user is required and cannot be empty")
	}
	if strings.TrimSpace(a.Item) == "" {
		return errors.New("item is required and cannot be empty")
	}
	if strings.TrimSpace(a.Action) == "" {
		return errors.New("action is required and cannot be empty")
	}
	// Check for control characters
	if containsControlChars(a.User) || containsControlChars(a.Item) || containsControlChars(a.Action) {
		return errors.New("user, item, and action cannot contain control characters")
	}
	return nil
}

// containsControlChars checks if string contains control characters
func containsControlChars(s string) bool {
	for _, r := range s {
		if r < 32 || r == 127 {
			return true
		}
	}
	return false
}
