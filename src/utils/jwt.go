package utils

import (
	"errors"
	"strings"

	"simple-sync/src/models"
	"simple-sync/src/services"
)

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Authorization header required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid authorization header format")
	}

	return parts[1], nil
}

// ValidateAndExtractClaims validates token and extracts claims
func ValidateAndExtractClaims(tokenString string, authService *services.AuthService) (*models.TokenClaims, error) {
	return authService.ValidateToken(tokenString)
}
