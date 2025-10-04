package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
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
