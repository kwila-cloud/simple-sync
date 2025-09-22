package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostEventsProtected(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth middleware
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Test data
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	// Test without Authorization header - should fail with 401
	req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Authorization header required", response["error"])
}

func TestPostEventsWithValidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Get valid token
	user, _ := h.AuthService().Authenticate("testuser", "testpass123")
	token, _ := h.AuthService().GenerateToken(user)

	// Test data - userUuid will be overridden by authenticated user
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user-123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	// Test with valid Authorization header
	req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 200 OK with events
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Should return the posted events
	assert.JSONEq(t, eventJSON, w.Body.String())
}

func TestPostEventsWithInvalidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Test data
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	// Test with invalid Authorization header
	req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
