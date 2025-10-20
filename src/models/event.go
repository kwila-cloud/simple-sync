package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	apperrors "simple-sync/src/errors"

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

// Validate performs validation on the Event struct
func (e *Event) Validate() error {
	if e.UUID == "" {
		return apperrors.ErrUuidRequired
	}

	// Validate UUID format and timestamp matching
	parsedUuid, err := uuid.Parse(e.UUID)
	if err != nil {
		return apperrors.ErrInvalidUuidFormat
	}

	// Validate that timestamp matches the UUID v7 timestamp
	timestamp, _ := parsedUuid.Time().UnixTime()
	if uint64(timestamp) != e.Timestamp {
		return apperrors.ErrInvalidTimestamp
	}

	// Do not allow timestamps of 0
	if uint64(timestamp) == 0 {
		return apperrors.ErrInvalidTimestamp
	}

	// Maximum timestamp: Allow up to 24 hours in the future for clock skew tolerance
	now := time.Now().Unix()
	maxTimestamp := now + (24 * 60 * 60) // 24 hours from now
	if int64(timestamp) > maxTimestamp {
		return apperrors.ErrInvalidTimestamp
	}

	if e.User == "" {
		return apperrors.ErrUserRequired
	}

	if e.Item == "" {
		return apperrors.ErrItemRequired
	}

	if e.Action == "" {
		return apperrors.ErrActionRequired
	}

	return nil
}

// ToAclRule converts an ACL event to AclRule
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
