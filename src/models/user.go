package models

import (
	"time"

	apperrors "simple-sync/src/errors"
)

// User represents an authenticated user in the system
type User struct {
	Id        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate performs validation on the User struct
func (u *User) Validate() error {
	if u.Id == "" {
		return apperrors.ErrIdRequired
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	return nil
}

// NewUser creates a new User with validation
func NewUser(id string) (*User, error) {
	user := &User{
		Id:        id,
		CreatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}
