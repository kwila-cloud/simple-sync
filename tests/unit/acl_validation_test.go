package unit

import (
	"testing"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

func TestACLEventValidate(t *testing.T) {
	tests := []struct {
		name     string
		event    models.ACLEvent
		expected error
	}{
		{
			name: "valid ACL event allow",
			event: models.ACLEvent{
				User:   "user123",
				Item:   "item456",
				Action: "read",
				Type:   "allow",
			},
			expected: nil,
		},
		{
			name: "valid ACL event deny",
			event: models.ACLEvent{
				User:   "user123",
				Item:   "item456",
				Action: "read",
				Type:   "deny",
			},
			expected: nil,
		},
		{
			name: "empty user",
			event: models.ACLEvent{
				User:   "",
				Item:   "item456",
				Action: "read",
				Type:   "allow",
			},
			expected: assert.AnError, // Will check error message
		},
		{
			name: "empty item",
			event: models.ACLEvent{
				User:   "user123",
				Item:   "",
				Action: "read",
				Type:   "allow",
			},
			expected: assert.AnError,
		},
		{
			name: "empty action",
			event: models.ACLEvent{
				User:   "user123",
				Item:   "item456",
				Action: "",
				Type:   "allow",
			},
			expected: assert.AnError,
		},
		{
			name: "invalid type",
			event: models.ACLEvent{
				User:   "user123",
				Item:   "item456",
				Action: "read",
				Type:   "invalid",
			},
			expected: assert.AnError,
		},
		{
			name: "user with control character",
			event: models.ACLEvent{
				User:   "user\x00",
				Item:   "item456",
				Action: "read",
				Type:   "allow",
			},
			expected: assert.AnError,
		},
		{
			name: "item with control character",
			event: models.ACLEvent{
				User:   "user123",
				Item:   "item\x01",
				Action: "read",
				Type:   "allow",
			},
			expected: assert.AnError,
		},
		{
			name: "action with control character",
			event: models.ACLEvent{
				User:   "user123",
				Item:   "item456",
				Action: "read\x02",
				Type:   "allow",
			},
			expected: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
