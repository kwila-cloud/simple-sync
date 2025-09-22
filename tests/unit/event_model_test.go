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

	assert.Equal(t, event.UUID, "test-uuid")
	assert.Equal(t, event.Timestamp, uint64(1234567890))
	assert.Equal(t, event.UserUUID, "user-uuid")
	assert.Equal(t, event.ItemUUID, "item-uuid")
	assert.Equal(t, event.Action, "update")
	assert.Equal(t, event.Payload, `{"key": "value"}`)
}

func TestNewHealthCheckResponse(t *testing.T) {
	status := "healthy"
	version := "1.0.0"
	uptime := int64(123)

	response := models.NewHealthCheckResponse(status, version, uptime)

	if response.Status != status {
		t.Errorf("Expected status %s, got %s", status, response.Status)
	}

	if response.Version != version {
		t.Errorf("Expected version %s, got %s", version, response.Version)
	}

	if response.Uptime != uptime {
		t.Errorf("Expected uptime %d, got %d", uptime, response.Uptime)
	}
}
