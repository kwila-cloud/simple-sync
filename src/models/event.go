package models

import (
	"encoding/json"
	"fmt"
)

// Event represents a timestamped event in the system
type Event struct {
	UUID      string `json:"uuid"`
	Timestamp uint64 `json:"timestamp"`
	User      string `json:"user"`
	Item      string `json:"item"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
}

// AclRule represents an access control rule
type AclRule struct {
	User      string `json:"user"`
	Item      string `json:"item"`
	Action    string `json:"action"`
	Type      string `json:"type"`
	Timestamp uint64 `json:"timestamp"`
}

// IsAclEvent checks if the event is an ACL rule event
func (e *Event) IsAclEvent() bool {
	return e.Item == ".acl"
}

// ToAclRule converts an ACL event to ACLRule
func (e *Event) ToAclRule() (*AclRule, error) {
	if !e.IsAclEvent() {
		return nil, fmt.Errorf("not an ACL event")
	}
	var rule AclRule
	err := json.Unmarshal([]byte(e.Payload), &rule)
	if err != nil {
		return nil, err
	}
	if e.Action == ".acl.allow" {
		rule.Type = "allow"
	} else if e.Action == ".acl.deny" {
		rule.Type = "deny"
	} else {
		return nil, fmt.Errorf("invalid ACL action: %s", e.Action)
	}
	rule.Timestamp = e.Timestamp
	return &rule, nil
}
