package middleware

import (
	"net/http"

	"simple-sync/src/services"
	"simple-sync/src/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates API key authentication middleware
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from Authorization header
		authHeader := c.GetHeader("Authorization")
		apiKey, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Validate API key
		userID, err := authService.ValidateApiKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", userID)

		c.Next()
	}
}
