package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"

	"github.com/stretchr/testify/assert"
)

func TestApiKeyValidate(t *testing.T) {
	tests := []struct {
		name    string
		key     *models.ApiKey
		wantErr bool
	}{
		{
			name: "valid API key",
			key: models.NewApiKey(
				"user123",
				"hash123",
				"test key",
			),
			wantErr: false,
		},
		{
			name: "empty UUID should fail",
			key: &models.ApiKey{
				UUID:        "",
				User:        "user123",
				KeyHash:     "hash123",
				CreatedAt:   time.Now(),
				Description: "test key",
			},
			wantErr: true,
		},
		{
			name: "invalid UUID format should fail",
			key: &models.ApiKey{
				UUID:        "invalid-uuid",
				User:        "user123",
				KeyHash:     "hash123",
				CreatedAt:   time.Now(),
				Description: "test key",
			},
			wantErr: true,
		},
		{
			name: "empty user ID should fail",
			key: &models.ApiKey{
				UUID:        "550e8400-e29b-41d4-a716-446655440000",
				User:        "",
				KeyHash:     "hash123",
				CreatedAt:   time.Now(),
				Description: "test key",
			},
			wantErr: true,
		},
		{
			name: "empty key hash should fail",
			key: &models.ApiKey{
				UUID:        "550e8400-e29b-41d4-a716-446655440000",
				User:        "user123",
				KeyHash:     "",
				CreatedAt:   time.Now(),
				Description: "test key",
			},
			wantErr: true,
		},
		{
			name: "zero created at should fail",
			key: &models.ApiKey{
				UUID:        "550e8400-e29b-41d4-a716-446655440000",
				User:        "user123",
				KeyHash:     "hash123",
				CreatedAt:   time.Time{},
				Description: "test key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.key.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApiKeyUpdateLastUsed(t *testing.T) {
	key := models.NewApiKey("user123", "hash123", "test key")

	// Initially LastUsedAt should be nil
	assert.Nil(t, key.LastUsedAt)

	// Update last used
	key.UpdateLastUsed()

	// LastUsedAt should now be set
	assert.NotNil(t, key.LastUsedAt)
	assert.WithinDuration(t, time.Now(), *key.LastUsedAt, time.Second)
}
