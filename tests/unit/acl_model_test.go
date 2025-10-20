package unit

import (
	"testing"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

func TestAclRuleValidate(t *testing.T) {
	tests := []struct {
		name    string
		rule    *models.AclRule
		wantErr bool
	}{
		{
			name: "valid rule",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "item456",
				Action: "read",
				Type:   "allow",
			},
			wantErr: false,
		},
		{
			name: "empty user should fail",
			rule: &models.AclRule{
				User:   "",
				Item:   "item456",
				Action: "read",
				Type:   "allow",
			},
			wantErr: true,
		},
		{
			name: "empty item should fail",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "",
				Action: "read",
				Type:   "allow",
			},
			wantErr: true,
		},
		{
			name: "empty action should fail",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "item456",
				Action: "",
				Type:   "allow",
			},
			wantErr: true,
		},
		{
			name: "empty type should fail",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "item456",
				Action: "read",
				Type:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
