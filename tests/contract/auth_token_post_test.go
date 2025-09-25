package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostAuthToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/auth/token", h.PostAuthToken)

	// Test data from contract
	authRequest := map[string]string{
		"username": "testuser",
		"password": "testpass123",
	}
	requestBody, _ := json.Marshal(authRequest)

	// Create test request
	req, _ := http.NewRequest("POST", "/api/v1/auth/token", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response according to contract
	// Expected: 200 with token
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

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/auth/token", h.PostAuthToken)

	// Invalid request - missing password
	invalidRequest := map[string]string{
		"username": "testuser",
	}
	requestBody, _ := json.Marshal(invalidRequest)

	req, _ := http.NewRequest("POST", "/api/v1/auth/token", bytes.NewBuffer(requestBody))
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

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/auth/token", h.PostAuthToken)

	// Invalid credentials
	authRequest := map[string]string{
		"username": "wronguser",
		"password": "wrongpass",
	}
	requestBody, _ := json.Marshal(authRequest)

	req, _ := http.NewRequest("POST", "/api/v1/auth/token", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPostUserResetKey(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/user/resetKey", h.PostUserResetKey)

	// Test data - reset key for user
	resetRequest := map[string]string{
		"user": "testuser",
	}
	requestBody, _ := json.Marshal(resetRequest)

	// Create test request with API key auth
	req, _ := http.NewRequest("POST", "/api/v1/user/resetKey?user=testuser", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk_testkey123456789012345678901234567890")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response according to contract
	// Expected: 200 with token
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "expiresAt")
	assert.NotEmpty(t, response["token"])
	assert.NotEmpty(t, response["expiresAt"])
}

func TestPostUserGenerateToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/user/generateToken", h.PostUserGenerateToken)

	// Test data - generate token for user
	generateRequest := map[string]string{
		"user": "testuser",
	}
	requestBody, _ := json.Marshal(generateRequest)

	// Create test request with API key auth
	req, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=testuser", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk_testkey123456789012345678901234567890")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response according to contract
	// Expected: 200 with token
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "expiresAt")
	assert.NotEmpty(t, response["token"])
	assert.NotEmpty(t, response["expiresAt"])
}

func TestPostSetupExchangeToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Test data - exchange setup token
	exchangeRequest := map[string]interface{}{
		"token":       "ABCD-1234",
		"description": "Test Client",
	}
	requestBody, _ := json.Marshal(exchangeRequest)

	// Create test request (no auth header needed for exchange)
	req, _ := http.NewRequest("POST", "/api/v1/setup/exchangeToken", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response according to contract
	// Expected: 200 with API key
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "keyUuid")
	assert.Contains(t, response, "apiKey")
	assert.Contains(t, response, "user")
	assert.NotEmpty(t, response["keyUuid"])
	assert.NotEmpty(t, response["apiKey"])
	assert.NotEmpty(t, response["user"])
}
