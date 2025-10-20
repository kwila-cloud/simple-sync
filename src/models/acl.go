package models

import (
	"errors"
)

// AclRule represents an access control rule
type AclRule struct {
	User   string `json:"user" db:"user"`
	Item   string `json:"item" db:"item"`
	Action string `json:"action" db:"action"`
	Type   string `json:"type" db:"type"`
}

// Validate performs validation on the AclRule struct
func (r *AclRule) Validate() error {
	if r.User == "" {
		return errors.New("user is required")
	}

	if r.Item == "" {
		return errors.New("item is required")
	}

	if r.Action == "" {
		return errors.New("action is required")
	}

	if r.Type == "" {
		return errors.New("type is required")
	}

	return nil
}
