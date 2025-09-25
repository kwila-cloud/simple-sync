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

func TestSuccessfulAuthenticationFlow(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/user/generateToken", h.PostUserGenerateToken)
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Protected routes with auth middleware
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Step 1: Generate setup token
	setupReq, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=user-123", nil)
	setupReq.Header.Set("Authorization", "Bearer sk_testkey123456789012345678901234567890") // Use test key
	setupW := httptest.NewRecorder()

	router.ServeHTTP(setupW, setupReq)

	assert.Equal(t, http.StatusOK, setupW.Code)

	var setupResponse map[string]string
	err := json.Unmarshal(setupW.Body.Bytes(), &setupResponse)
	assert.NoError(t, err)
	setupToken := setupResponse["token"]
	assert.NotEmpty(t, setupToken)

	// Step 2: Exchange for API key
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
	assert.NotEmpty(t, apiKey)

	// Step 3: Use API key to access protected GET /events
	getReq, _ := http.NewRequest("GET", "/api/v1/events", nil)
	getReq.Header.Set("Authorization", "Bearer "+apiKey)
	getW := httptest.NewRecorder()

	router.ServeHTTP(getW, getReq)

	assert.Equal(t, http.StatusOK, getW.Code)

	// Step 4: Use API key to POST events
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", "Bearer "+apiKey)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should succeed
	assert.Equal(t, http.StatusOK, postW.Code)

	// Expected response with authenticated user UUID
	expectedJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user-123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`
	assert.JSONEq(t, expectedJSON, postW.Body.String())
}
