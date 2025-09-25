package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

func TestAPIKeyModelValidation(t *testing.T) {
	// Test valid API key
	validKey := &models.APIKey{
		UUID:         "550e8400-e29b-41d4-a716-446655440000",
		UserID:       "user-123",
		EncryptedKey: "encrypted-data",
		KeyHash:      "hash-data",
		CreatedAt:    time.Now(),
		Description:  "Test Key",
	}
	err := validKey.Validate()
	assert.NoError(t, err)

	// Test invalid API key - missing UUID
	invalidKey := &models.APIKey{
		UserID:       "user-123",
		EncryptedKey: "encrypted-data",
		KeyHash:      "hash-data",
		CreatedAt:    time.Now(),
	}
	err = invalidKey.Validate()
	assert.Error(t, err)
}

func TestSetupTokenModelValidation(t *testing.T) {
	// Test valid setup token
	validToken := &models.SetupToken{
		Token:     "ABCD-1234",
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}
	err := validToken.Validate()
	assert.NoError(t, err)
	assert.True(t, validToken.IsValid())

	// Test invalid token format
	invalidToken := &models.SetupToken{
		Token:     "INVALID-FORMAT",
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}
	err = invalidToken.Validate()
	assert.Error(t, err)
}
