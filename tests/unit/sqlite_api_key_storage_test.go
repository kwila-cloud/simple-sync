package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"
)

// helper: create initialized in-memory sqlite storage
func newTestSQLiteStorage(t *testing.T) *storage.SQLiteStorage {
	s := storage.NewSQLiteStorage()
	if err := s.Initialize("file::memory:?cache=shared"); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	return s
}

func TestCreateAndGetApiKeyByHash(t *testing.T) {
	s := newTestSQLiteStorage(t)
	defer s.Close()

	// create user expected by foreign key
	if err := s.AddUser(&models.User{Id: "user-x", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	key := models.NewApiKey("user-x", "hash123", "test")
	if err := s.AddApiKey(key); err != nil {
		t.Fatalf("AddApiKey failed: %v", err)
	}

	got, err := s.GetApiKeyByHash("hash123")
	if err != nil {
		t.Fatalf("GetApiKeyByHash failed: %v", err)
	}
	if got.UUID != key.UUID {
		t.Fatalf("expected uuid %s, got %s", key.UUID, got.UUID)
	}
}

func TestGetAllApiKeysAndUpdateInvalidate(t *testing.T) {
	s := newTestSQLiteStorage(t)
	defer s.Close()

	if err := s.AddUser(&models.User{Id: "user-y", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	k1 := models.NewApiKey("user-y", "hash-a", "one")
	k2 := models.NewApiKey("user-y", "hash-b", "two")
	if err := s.AddApiKey(k1); err != nil {
		t.Fatalf("AddApiKey k1 failed: %v", err)
	}
	if err := s.AddApiKey(k2); err != nil {
		t.Fatalf("AddApiKey k2 failed: %v", err)
	}

	all, err := s.GetAllApiKeys()
	if err != nil {
		t.Fatalf("GetAllApiKeys failed: %v", err)
	}
	if len(all) < 2 {
		t.Fatalf("expected at least 2 keys, got %d", len(all))
	}

	// update last used
	k1.UpdateLastUsed()
	if err := s.UpdateApiKey(k1); err != nil {
		t.Fatalf("UpdateApiKey failed: %v", err)
	}

	// invalidate user keys
	if err := s.InvalidateUserApiKeys("user-y"); err != nil {
		t.Fatalf("InvalidateUserApiKeys failed: %v", err)
	}

	allAfter, err := s.GetAllApiKeys()
	if err != nil {
		t.Fatalf("GetAllApiKeys failed: %v", err)
	}
	for _, k := range allAfter {
		if k.User == "user-y" {
			t.Fatalf("expected keys for user-y to be deleted")
		}
	}
}
