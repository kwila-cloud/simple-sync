package unit

import (
	"testing"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"
)

func TestCreateGetUpdateInvalidateSetupToken(t *testing.T) {
	s := storage.NewSQLiteStorage()
	if err := s.Initialize("file::memory:?cache=shared"); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer s.Close()

	// create user
	if err := s.AddUser(&models.User{Id: "user-z", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	st := models.NewSetupToken("ABCD-1234", "user-z", time.Now().Add(time.Hour))
	if err := s.CreateSetupToken(st); err != nil {
		t.Fatalf("CreateSetupToken failed: %v", err)
	}

	got, err := s.GetSetupToken("ABCD-1234")
	if err != nil {
		t.Fatalf("GetSetupToken failed: %v", err)
	}
	if got.User != "user-z" {
		t.Fatalf("expected user %s, got %s", "user-z", got.User)
	}

	// mark used and update
	got.MarkUsed()
	if err := s.UpdateSetupToken(got); err != nil {
		t.Fatalf("UpdateSetupToken failed: %v", err)
	}

	// invalidate tokens for user (should mark used_at for all their tokens)
	if err := s.InvalidateUserSetupTokens("user-z"); err != nil {
		t.Fatalf("InvalidateUserSetupTokens failed: %v", err)
	}

	// fetch again and ensure used_at is set
	got2, err := s.GetSetupToken("ABCD-1234")
	if err != nil {
		t.Fatalf("GetSetupToken failed: %v", err)
	}
	if got2.IsValid() {
		t.Fatalf("expected token to be marked used")
	}
}
