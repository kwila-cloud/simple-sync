package performance

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpointPerformance(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Use the default test API key instead of generating a new one
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Test auth endpoint performance
	start := time.Now()
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert performance (<100ms)
	assert.Less(t, duration, 100*time.Millisecond, "Auth endpoint should respond in less than 100ms")
}

func TestProtectedEndpointPerformance(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Use the default test API key instead of generating a new one
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Test protected endpoint performance
	start := time.Now()
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert performance (<100ms)
	assert.Less(t, duration, 100*time.Millisecond, "Protected endpoint should respond in less than 100ms")
}
