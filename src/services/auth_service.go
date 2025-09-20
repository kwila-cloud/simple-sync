package services

import (
	"errors"
	"time"

	"simple-sync/src/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	jwtSecret []byte
	users     map[string]*models.User // In-memory user store for MVP
}

// NewAuthService creates a new auth service
func NewAuthService(jwtSecret string) *AuthService {
	service := &AuthService{
		jwtSecret: []byte(jwtSecret),
		users:     make(map[string]*models.User),
	}

	// Add default user for MVP with password hashing
	defaultUser, _ := models.NewUserWithPassword("user-123", "testuser", "testpass123", false)
	service.users[defaultUser.Username] = defaultUser

	return service
}

// Authenticate validates user credentials and returns user if valid
func (s *AuthService) Authenticate(username, password string) (*models.User, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, errors.New("Invalid username or password")
	}

	// Verify password using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
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
