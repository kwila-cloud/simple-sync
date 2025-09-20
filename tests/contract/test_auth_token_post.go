package contract

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

func TestPostAuthToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret")

	// Register routes
	router.POST("/auth/token", h.PostAuthToken)

	// Test data from contract
	authRequest := map[string]string{
		"username": "testuser",
		"password": "testpass123",
	}
	requestBody, _ := json.Marshal(authRequest)

	// Create test request
	req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response according to contract
	// Expected: 200 with token, but will fail since endpoint doesn't exist
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])
}

func TestPostAuthTokenInvalidRequest(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret")

	// Register routes
	router.POST("/auth/token", h.PostAuthToken)

	// Invalid request - missing password
	invalidRequest := map[string]string{
		"username": "testuser",
	}
	requestBody, _ := json.Marshal(invalidRequest)

	req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostAuthTokenInvalidCredentials(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret")

	// Register routes
	router.POST("/auth/token", h.PostAuthToken)

	// Invalid credentials
	authRequest := map[string]string{
		"username": "wronguser",
		"password": "wrongpass",
	}
	requestBody, _ := json.Marshal(authRequest)

	req, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
