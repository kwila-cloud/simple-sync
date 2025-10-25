package middleware

import (
	"net/http"

	"simple-sync/src/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates API key authentication middleware
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from X-API-Key header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Spec tests expect a 406 when a required header is missing
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "X-API-Key header required"})
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
