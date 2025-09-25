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

func TestPostUserResetKey(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/user/resetKey", h.PostUserResetKey)

	// Test data - reset key for user
	resetRequest := map[string]string{
		"user": "testuser",
	}
	requestBody, _ := json.Marshal(resetRequest)

	// Create test request with API key auth
	req, _ := http.NewRequest("POST", "/api/v1/user/resetKey?user=user-123", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk_testkey123456789012345678901234567890")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response according to contract
	// Expected: 200 with success message
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Equal(t, "API keys invalidated successfully", response["message"])
}

func TestPostUserGenerateToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/user/generateToken", h.PostUserGenerateToken)

	// Test data - generate token for user
	generateRequest := map[string]string{
		"user": "testuser",
	}
	requestBody, _ := json.Marshal(generateRequest)

	// Create test request with API key auth
	req, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=user-123", bytes.NewBuffer(requestBody))
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

	// Generate setup token first
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Test data - exchange setup token
	exchangeRequest := map[string]interface{}{
		"token":       setupToken.Token,
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
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "keyUuid")
	assert.Contains(t, response, "apiKey")
	assert.Contains(t, response, "user")
	assert.NotEmpty(t, response["keyUuid"])
	assert.NotEmpty(t, response["apiKey"])
	assert.NotEmpty(t, response["user"])
}
