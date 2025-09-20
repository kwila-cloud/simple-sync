package unit

import (
	"encoding/json"
	"testing"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

func TestEventJSONMarshaling(t *testing.T) {
	event := models.Event{
		UUID:      "123e4567-e89b-12d3-a456-426614174000",
		Timestamp: 1640995200,
		UserUUID:  "user123",
		ItemUUID:  "item456",
		Action:    "create",
		Payload:   "{}",
	}

	// Test marshaling
	data, err := json.Marshal(event)
	assert.NoError(t, err)

	// Test unmarshaling
	var unmarshaled models.Event
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, event, unmarshaled)
}

func TestEventFields(t *testing.T) {
	event := models.Event{
		UUID:      "test-uuid",
		Timestamp: 1234567890,
		UserUUID:  "user-uuid",
		ItemUUID:  "item-uuid",
		Action:    "update",
		Payload:   `{"key": "value"}`,
	}

	assert.Equal(t, "test-uuid", event.UUID)
	assert.Equal(t, uint64(1234567890), event.Timestamp)
	assert.Equal(t, "user-uuid", event.UserUUID)
	assert.Equal(t, "item-uuid", event.ItemUUID)
	assert.Equal(t, "update", event.Action)
	assert.Equal(t, `{"key": "value"}`, event.Payload)
}