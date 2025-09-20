package models

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents an authenticated user in the system
type User struct {
	UUID         string    `json:"uuid"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	IsAdmin      bool      `json:"is_admin"`
}

// Validate performs validation on the User struct
func (u *User) Validate() error {
	if u.UUID == "" {
		return errors.New("uuid is required")
	}

	if u.Username == "" {
		return errors.New("username is required")
	}

	// Username validation: 3-50 characters, alphanumeric + underscore/hyphen
	if len(u.Username) < 3 || len(u.Username) > 50 {
		return errors.New("username must be 3-50 characters")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(u.Username) {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}

	if u.PasswordHash == "" {
		return errors.New("password hash is required")
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	return nil
}

// NewUser creates a new User with validation
func NewUser(uuid, username, passwordHash string, isAdmin bool) (*User, error) {
	user := &User{
		UUID:         uuid,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		IsAdmin:      isAdmin,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// NewUserWithPassword creates a new User with password hashing
func NewUserWithPassword(uuid, username, plainPassword string, isAdmin bool) (*User, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return NewUser(uuid, username, string(hashedPassword), isAdmin)
}
