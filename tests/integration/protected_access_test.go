package integration

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

func TestProtectedEndpointAccess(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	router.POST("/auth/token", h.PostAuthToken)

	// Protected routes with auth middleware
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Test 1: Access GET /events without token - should fail
	getReq, _ := http.NewRequest("GET", "/events", nil)
	getW := httptest.NewRecorder()

	router.ServeHTTP(getW, getReq)

	// Expected: 401 (will fail until middleware)
	assert.Equal(t, http.StatusUnauthorized, getW.Code)

	// Test 2: Access POST /events without token - should fail
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	postReq, _ := http.NewRequest("POST", "/events", bytes.NewBufferString(eventJSON))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Expected: 401
	assert.Equal(t, http.StatusUnauthorized, postW.Code)

	// Test 3: Get token, then access with token - should succeed
	authRequest := map[string]string{
		"username": "testuser",
		"password": "testpass123",
	}
	authBody, _ := json.Marshal(authRequest)

	authReq, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(authBody))
	authReq.Header.Set("Content-Type", "application/json")
	authW := httptest.NewRecorder()

	router.ServeHTTP(authW, authReq)

	assert.Equal(t, http.StatusOK, authW.Code)

	var authResponse map[string]string
	err := json.Unmarshal(authW.Body.Bytes(), &authResponse)
	assert.NoError(t, err)
	token := authResponse["token"]

	// Now access with token
	authGetReq, _ := http.NewRequest("GET", "/events", nil)
	authGetReq.Header.Set("Authorization", "Bearer "+token)
	authGetW := httptest.NewRecorder()

	router.ServeHTTP(authGetW, authGetReq)

	assert.Equal(t, http.StatusOK, authGetW.Code)
}
