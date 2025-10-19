package handlers

import (
	"simple-sync/src/models"
	"testing"
)

func TestValidateAclRuleSpecificErrors(t *testing.T) {
	tests := []struct {
		name        string
		rule        models.AclRule
		expectedErr error
	}{
		{
			name: "empty user pattern",
			rule: models.AclRule{
				User:   "",
				Item:   "item1",
				Action: "read",
				Type:   "allow",
			},
			expectedErr: ErrAclUserEmpty,
		},
		{
			name: "empty item pattern",
			rule: models.AclRule{
				User:   "user1",
				Item:   "",
				Action: "read",
				Type:   "allow",
			},
			expectedErr: ErrAclItemEmpty,
		},
		{
			name: "empty action pattern",
			rule: models.AclRule{
				User:   "user1",
				Item:   "item1",
				Action: "",
				Type:   "allow",
			},
			expectedErr: ErrAclActionEmpty,
		},
		{
			name: "multiple wildcards in user",
			rule: models.AclRule{
				User:   "user*test*",
				Item:   "item1",
				Action: "read",
				Type:   "allow",
			},
			expectedErr: ErrAclUserMultipleWildcards,
		},
		{
			name: "multiple wildcards in item",
			rule: models.AclRule{
				User:   "user1",
				Item:   "item*test*",
				Action: "read",
				Type:   "allow",
			},
			expectedErr: ErrAclItemMultipleWildcards,
		},
		{
			name: "multiple wildcards in action",
			rule: models.AclRule{
				User:   "user1",
				Item:   "item1",
				Action: "read*write*",
				Type:   "allow",
			},
			expectedErr: ErrAclActionMultipleWildcards,
		},
		{
			name: "invalid ACL type",
			rule: models.AclRule{
				User:   "user1",
				Item:   "item1",
				Action: "read",
				Type:   "invalid",
			},
			expectedErr: ErrInvalidAclType,
		},
		{
			name: "valid rule",
			rule: models.AclRule{
				User:   "user1",
				Item:   "item1",
				Action: "read",
				Type:   "allow",
			},
			expectedErr: nil,
		},
		{
			name: "valid wildcard at end",
			rule: models.AclRule{
				User:   "user*",
				Item:   "item*",
				Action: "read*",
				Type:   "allow",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAclRule(tt.rule)
			if err != tt.expectedErr {
				t.Errorf("validateAclRule() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
