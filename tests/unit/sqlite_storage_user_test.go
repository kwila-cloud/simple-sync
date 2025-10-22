package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/stretchr/testify/assert"
)

func TestSQLiteStorage_AddUserAndGetUserById(t *testing.T) {
	// Use in-memory SQLite for tests
	s := storage.NewSQLiteStorage()
	if err := s.Initialize(":memory:"); err != nil {
		t.Fatalf("failed to initialize sqlite storage: %v", err)
	}
	defer s.Close()

	user, err := models.NewUser("test-user-1")
	assert.NoError(t, err)

	// Add user
	err = s.AddUser(user)
	assert.NoError(t, err)

	// Retrieve existing user
	got, err := s.GetUserById("test-user-1")
	assert.NoError(t, err)
	assert.Equal(t, user.Id, got.Id)
	// Allow some tolerance for created_at since NewUser uses time.Now()
	assert.WithinDuration(t, user.CreatedAt, got.CreatedAt, 2*time.Second)

	// Adding duplicate should return ErrDuplicateKey
	err = s.AddUser(user)
	assert.Error(t, err)
	assert.Equal(t, storage.ErrDuplicateKey, err)

	// Non-existent user returns ErrNotFound
	_, err = s.GetUserById("no-such-user")
	assert.Error(t, err)
	assert.Equal(t, storage.ErrNotFound, err)
}
