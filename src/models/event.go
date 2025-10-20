package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Event represents a timestamped event in the system
type Event struct {
	UUID      string `json:"uuid" db:"uuid"`
	Timestamp uint64 `json:"timestamp" db:"timestamp"`
	User      string `json:"user" db:"user"`
	Item      string `json:"item" db:"item"`
	Action    string `json:"action" db:"action"`
	Payload   string `json:"payload" db:"payload"`
}

func NewEvent(User, Item, Action, Payload string) *Event {
	eventUuid, _ := uuid.NewV7()
	unixTimeSeconds, _ := eventUuid.Time().UnixTime()

	return &Event{
		UUID:      eventUuid.String(),
		Timestamp: uint64(unixTimeSeconds),
		User:      User,
		Item:      Item,
		Action:    Action,
		Payload:   Payload,
	}
}

func (e *Event) IsApiOnlyEvent() bool {
	// .user.create is the ONLY internal event action that can be triggered
	// by a user
	return e.Action != ".user.create" && strings.HasPrefix(e.Action, ".")
}

func (e *Event) IsAclEvent() bool {
	return e.Item == ".acl"
}

// ToAclRule converts an ACL event to AclRule
// Validate performs validation on the Event struct
func (e *Event) Validate() error {
	if e.UUID == "" {
		return errors.New("UUID is required")
	}

	// Validate UUID format and timestamp matching
	parsedUuid, err := uuid.Parse(e.UUID)
	if err != nil {
		return errors.New("UUID must be valid format")
	}

	// For v7 UUIDs, validate that timestamp matches
	timestamp, _ := parsedUuid.Time().UnixTime()
	if uint64(timestamp) != e.Timestamp {
		return errors.New("UUID timestamp must match event timestamp")
	}

	if e.User == "" {
		return errors.New("user is required")
	}

	if e.Item == "" {
		return errors.New("item is required")
	}

	if e.Action == "" {
		return errors.New("action is required")
	}

	return nil
}

func (e *Event) ToAclRule() (*AclRule, error) {
	if !e.IsAclEvent() {
		return nil, fmt.Errorf("not an ACL event")
	}
	var rule AclRule
	err := json.Unmarshal([]byte(e.Payload), &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}
