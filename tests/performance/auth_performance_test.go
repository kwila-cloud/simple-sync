package performance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpointPerformance(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret", "test")

	// Register routes
	router.POST("/auth/token", h.PostAuthToken)

	// Test auth endpoint performance
	authRequest := map[string]string{
		"username": "testuser",
		"password": "testpass123",
	}
	authBody, _ := json.Marshal(authRequest)

	start := time.Now()
	req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
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

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret", "test")

	// Register routes with auth
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Get token
	user, _ := h.AuthService().Authenticate("testuser", "testpass123")
	token, _ := h.AuthService().GenerateToken(user)

	// Test protected endpoint performance
	start := time.Now()
	req, _ := http.NewRequest("GET", "/events", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert performance (<100ms)
	assert.Less(t, duration, 100*time.Millisecond, "Protected endpoint should respond in less than 100ms")
}
