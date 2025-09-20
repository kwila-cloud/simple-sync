package services

import (
	"errors"
	"time"

	"simple-sync/src/models"

	"github.com/golang-jwt/jwt/v5"
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

	// Add default user for MVP
	defaultUser, _ := models.NewUser("user-123", "testuser", "testpass123", false)
	service.users[defaultUser.Username] = defaultUser

	return service
}

// Authenticate validates user credentials and returns user if valid
func (s *AuthService) Authenticate(username, password string) (*models.User, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, errors.New("invalid username or password")
	}

	// For MVP, simple password comparison (should use bcrypt in production)
	if user.PasswordHash != password {
		return nil, errors.New("invalid username or password")
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

	return nil, errors.New("invalid token")
}

// GetUserByUUID retrieves a user by UUID
func (s *AuthService) GetUserByUUID(uuid string) (*models.User, error) {
	for _, user := range s.users {
		if user.UUID == uuid {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
