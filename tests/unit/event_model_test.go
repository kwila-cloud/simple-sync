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
