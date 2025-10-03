package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestXAPIKeyHeaderAccepted tests that requests with X-API-Key header are accepted
func TestXAPIKeyHeaderAccepted(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register protected route
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Use default test API key
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Create test request with X-API-Key header
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("X-API-Key", plainKey)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert successful authentication
	assert.Equal(t, http.StatusOK, w.Code)
}
