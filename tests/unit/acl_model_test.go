package unit

import (
	"testing"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

func TestAclRuleValidate(t *testing.T) {
	tests := []struct {
		name        string
		rule        *models.AclRule
		wantErr     bool
		expectedErr string
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
			wantErr:     true,
			expectedErr: "user pattern cannot be empty",
		},
		{
			name: "empty item should fail",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "",
				Action: "read",
				Type:   "allow",
			},
			wantErr:     true,
			expectedErr: "item pattern cannot be empty",
		},
		{
			name: "empty action should fail",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "item456",
				Action: "",
				Type:   "allow",
			},
			wantErr:     true,
			expectedErr: "action pattern cannot be empty",
		},
		{
			name: "invalid ACL type",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "item456",
				Action: "read",
				Type:   "invalid",
			},
			wantErr:     true,
			expectedErr: "type must be either 'allow' or 'deny'",
		},
		{
			name: "multiple wildcards in user",
			rule: &models.AclRule{
				User:   "user*test*",
				Item:   "item456",
				Action: "read",
				Type:   "allow",
			},
			wantErr:     true,
			expectedErr: "user pattern can have at most one wildcard at the end",
		},
		{
			name: "valid wildcard at end",
			rule: &models.AclRule{
				User:   "user*",
				Item:   "item*",
				Action: "read*",
				Type:   "allow",
			},
			wantErr: false,
		},
		{
			name: "wildcard in middle should fail",
			rule: &models.AclRule{
				User:   "user123",
				Item:   "item*test",
				Action: "read",
				Type:   "allow",
			},
			wantErr:     true,
			expectedErr: "item pattern can have at most one wildcard at the end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Equal(t, tt.expectedErr, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
