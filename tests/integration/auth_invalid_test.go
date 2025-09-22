package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInvalidCredentialsHandling(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret", "test")

	// Register routes
	router.POST("/auth/token", h.PostAuthToken)
	router.GET("/events", h.GetEvents)

	// Test 1: Invalid username/password
	authRequest := map[string]string{
		"username": "wronguser",
		"password": "wrongpass",
	}
	authBody, _ := json.Marshal(authRequest)

	authReq, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(authBody))
	authReq.Header.Set("Content-Type", "application/json")
	authW := httptest.NewRecorder()

	router.ServeHTTP(authW, authReq)

	// Expected: 401
	assert.Equal(t, http.StatusUnauthorized, authW.Code)

	var authResponse map[string]string
	err := json.Unmarshal(authW.Body.Bytes(), &authResponse)
	assert.NoError(t, err)
	assert.Contains(t, authResponse, "error")
	assert.Equal(t, "Invalid username or password", authResponse["error"])

	// Test 2: Malformed request (missing password)
	malformedRequest := map[string]string{
		"username": "testuser",
	}
	malformedBody, _ := json.Marshal(malformedRequest)

	malformedReq, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(malformedBody))
	malformedReq.Header.Set("Content-Type", "application/json")
	malformedW := httptest.NewRecorder()

	router.ServeHTTP(malformedW, malformedReq)

	// Expected: 400
	assert.Equal(t, http.StatusBadRequest, malformedW.Code)

	// Test 3: Access with invalid token
	getReq, _ := http.NewRequest("GET", "/events", nil)
	getReq.Header.Set("Authorization", "Bearer invalid-token")
	getW := httptest.NewRecorder()

	router.ServeHTTP(getW, getReq)

	// Expected: 401
	assert.Equal(t, http.StatusUnauthorized, getW.Code)

	// Test 4: Access with expired token (simulate)
	// Note: This would require implementing token expiration
	getExpiredReq, _ := http.NewRequest("GET", "/events", nil)
	getExpiredReq.Header.Set("Authorization", "Bearer expired-jwt-token")
	getExpiredW := httptest.NewRecorder()

	router.ServeHTTP(getExpiredW, getExpiredReq)

	// Expected: 401 (will fail until expiration check implemented)
	assert.Equal(t, http.StatusUnauthorized, getExpiredW.Code)
}
