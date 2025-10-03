package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProtectedEndpointAccess(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers with memory storage
	store := storage.NewMemoryStorage(nil)
	h := handlers.NewTestHandlersWithStorage(store)

	// Create root user and API key for authentication
	rootUser := &models.User{Id: ".root"}
	err := store.SaveUser(rootUser)
	assert.NoError(t, err)

	_, adminApiKey, err := h.AuthService().GenerateApiKey(".root", "Admin Key")
	assert.NoError(t, err)

	// Register routes
	v1 := router.Group("/api/v1")

	// Auth routes with middleware
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/user/generateToken", h.PostUserGenerateToken)

	// Setup routes (no middleware)
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Protected routes with auth middleware
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Test 1: Access GET /events without token - should fail
	getReq, _ := http.NewRequest("GET", "/api/v1/events", nil)
	getW := httptest.NewRecorder()

	router.ServeHTTP(getW, getReq)

	// Expected: 401 (will fail until middleware)
	assert.Equal(t, http.StatusUnauthorized, getW.Code)

	// Test 2: Access POST /events without token - should fail
	eventJSON := `[{
 		"uuid": "123e4567-e89b-12d3-a456-426614174000",
 		"timestamp": 1640995200,
 		"user": "user123",
 		"item": "item456",
 		"action": "create",
 		"payload": "{}"
 	}]`

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Expected: 401
	assert.Equal(t, http.StatusUnauthorized, postW.Code)

	// Test 3: Get API key, then access with API key - should succeed
	// Generate setup token
	setupReq, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=user-123", nil)
	setupReq.Header.Set("X-API-Key", adminApiKey)
	setupW := httptest.NewRecorder()

	router.ServeHTTP(setupW, setupReq)

	assert.Equal(t, http.StatusOK, setupW.Code)

	var setupResponse map[string]string
	err = json.Unmarshal(setupW.Body.Bytes(), &setupResponse)
	assert.NoError(t, err)
	setupToken := setupResponse["token"]

	// Exchange for API key
	exchangeRequest := map[string]interface{}{
		"token": setupToken,
	}
	exchangeBody, _ := json.Marshal(exchangeRequest)

	exchangeReq, _ := http.NewRequest("POST", "/api/v1/setup/exchangeToken", bytes.NewBuffer(exchangeBody))
	exchangeReq.Header.Set("Content-Type", "application/json")
	exchangeW := httptest.NewRecorder()

	router.ServeHTTP(exchangeW, exchangeReq)

	assert.Equal(t, http.StatusOK, exchangeW.Code)

	var exchangeResponse map[string]interface{}
	err = json.Unmarshal(exchangeW.Body.Bytes(), &exchangeResponse)
	assert.NoError(t, err)
	apiKey := exchangeResponse["apiKey"].(string)

	// Now access with API key
	authGetReq, _ := http.NewRequest("GET", "/api/v1/events", nil)
	authGetReq.Header.Set("X-API-Key", apiKey)
	authGetW := httptest.NewRecorder()

	router.ServeHTTP(authGetW, authGetReq)

	assert.Equal(t, http.StatusOK, authGetW.Code)
}
