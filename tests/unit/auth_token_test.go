package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApiKeyModelValidation(t *testing.T) {
	keyUuid, _ := uuid.NewV7()
	unixTimeSeconds, _ := keyUuid.Time().UnixTime()

	// Test valid API key
	validKey := &models.ApiKey{
		UUID:        keyUuid.String(),
		UserID:      storage.TestingUserId,
		KeyHash:     "hash-data",
		CreatedAt:   time.Unix(unixTimeSeconds, 0),
		Description: "Test Key",
	}
	err := validKey.Validate()
	assert.NoError(t, err)

	// Test invalid API key - missing UUID
	invalidKey := &models.ApiKey{
		UserID:    storage.TestingUserId,
		KeyHash:   "hash-data",
		CreatedAt: time.Now(),
	}
	err = invalidKey.Validate()
	assert.Error(t, err)
}

func TestSetupTokenModelValidation(t *testing.T) {
	// Test valid setup token
	validToken := &models.SetupToken{
		Token:     "ABCD-1234",
		UserID:    storage.TestingUserId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UsedAt:    time.Time{},
	}
	err := validToken.Validate()
	assert.NoError(t, err)
	assert.True(t, validToken.IsValid())

	// Test invalid token format
	invalidToken := &models.SetupToken{
		Token:     "INVALID-FORMAT",
		UserID:    storage.TestingUserId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UsedAt:    time.Time{},
	}
	err = invalidToken.Validate()
	assert.Error(t, err)
}
