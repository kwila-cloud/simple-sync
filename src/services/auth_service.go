package services

import (
	"errors"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	jwtSecret []byte
	storage   storage.Storage
}

// NewAuthService creates a new auth service
func NewAuthService(jwtSecret string, storage storage.Storage) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
		storage:   storage,
	}
}

// Authenticate validates user credentials and returns user if valid
func (s *AuthService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.storage.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("Invalid username or password")
	}

	// Verify password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid username or password")
	}

	return user, nil
}

// GenerateToken creates a JWT token for the authenticated user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(24 * time.Hour) // 24 hours

	claims := models.NewTokenClaims(user.UUID, user.Username, user.IsAdmin, issuedAt, expiresAt)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")
}
