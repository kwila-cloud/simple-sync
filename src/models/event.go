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

// ACLRule represents an access control rule
type ACLRule struct {
	User      string `json:"user"`
	Item      string `json:"item"`
	Action    string `json:"action"`
	Type      string `json:"type"` // "allow" or "deny"
	Timestamp uint64 `json:"timestamp"`
}

// IsACLEvent checks if the event is an ACL rule event
func (e *Event) IsACLEvent() bool {
	return e.Item == ".acl" && (e.Action == ".acl.allow" || e.Action == ".acl.deny")
}

// ToACLRule converts an ACL event to ACLRule
func (e *Event) ToACLRule() (*ACLRule, error) {
	if !e.IsACLEvent() {
		return nil, fmt.Errorf("not an ACL event")
	}
	var rule ACLRule
	err := json.Unmarshal([]byte(e.Payload), &rule)
	if err != nil {
		return nil, err
	}
	rule.Timestamp = e.Timestamp
	return &rule, nil
}
