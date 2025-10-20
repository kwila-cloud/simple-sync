package unit

import (
	"encoding/json"
	"testing"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

func TestEventJsonMarshaling(t *testing.T) {
	event := models.NewEvent(
		"user123",
		"item456",
		"create",
		"{}",
	)

	// Test marshaling
	data, err := json.Marshal(event)
	assert.NoError(t, err)

	// Test unmarshaling
	var unmarshaled models.Event
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, *event, unmarshaled)
}

func TestEventFields(t *testing.T) {
	event := models.NewEvent(
		"user-uuid",
		"item-uuid",
		"update",
		`{"key": "value"}`,
	)

	assert.Equal(t, event.User, "user-uuid")
	assert.Equal(t, event.Item, "item-uuid")
	assert.Equal(t, event.Action, "update")
	assert.Equal(t, event.Payload, `{"key": "value"}`)
}

func TestEventValidate(t *testing.T) {
	tests := []struct {
		name    string
		event   *models.Event
		wantErr bool
	}{
		{
			name: "valid event",
			event: models.NewEvent(
				"user123",
				"item456",
				"create",
				"{}",
			),
			wantErr: false,
		},
		{
			name: "empty UUID should fail",
			event: &models.Event{
				UUID:      "",
				Timestamp: 1234567890,
				User:      "user123",
				Item:      "item456",
				Action:    "create",
				Payload:   "{}",
			},
			wantErr: true,
		},
		{
			name: "invalid UUID format should fail",
			event: &models.Event{
				UUID:      "invalid-uuid",
				Timestamp: 1234567890,
				User:      "user123",
				Item:      "item456",
				Action:    "create",
				Payload:   "{}",
			},
			wantErr: true,
		},
		{
			name: "UUID timestamp mismatch should fail",
			event: &models.Event{
				UUID:      "550e8400-e29b-41d4-a716-446655440000",
				Timestamp: 1234567890, // Different timestamp than UUID
				User:      "user123",
				Item:      "item456",
				Action:    "create",
				Payload:   "{}",
			},
			wantErr: true,
		},
		{
			name: "empty user should fail",
			event: &models.Event{
				UUID:      "550e8400-e29b-41d4-a716-446655440000",
				Timestamp: 1234567890,
				User:      "",
				Item:      "item456",
				Action:    "create",
				Payload:   "{}",
			},
			wantErr: true,
		},
		{
			name: "empty item should fail",
			event: &models.Event{
				UUID:      "550e8400-e29b-41d4-a716-446655440000",
				Timestamp: 1234567890,
				User:      "user123",
				Item:      "",
				Action:    "create",
				Payload:   "{}",
			},
			wantErr: true,
		},
		{
			name: "empty action should fail",
			event: &models.Event{
				UUID:      "550e8400-e29b-41d4-a716-446655440000",
				Timestamp: 1234567890,
				User:      "user123",
				Item:      "item456",
				Action:    "",
				Payload:   "{}",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
