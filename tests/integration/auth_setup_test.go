package integration

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

func TestUserSetupFlowIntegration(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/user/generateToken", h.PostUserGenerateToken)
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Step 1: Generate setup token for user
	generateRequest := map[string]string{
		"user": "testuser",
	}
	requestBody, _ := json.Marshal(generateRequest)

	req1, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=testuser", bytes.NewBuffer(requestBody))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer sk_admin123456789012345678901234567890")
	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, req1)

	// Verify token generation
	assert.Equal(t, http.StatusOK, w1.Code)

	var generateResponse map[string]string
	err := json.Unmarshal(w1.Body.Bytes(), &generateResponse)
	assert.NoError(t, err)
	assert.Contains(t, generateResponse, "token")
	token := generateResponse["token"]

	// Step 2: Exchange setup token for API key
	exchangeRequest := map[string]interface{}{
		"token":       token,
		"description": "Integration Test Client",
	}
	exchangeBody, _ := json.Marshal(exchangeRequest)

	req2, _ := http.NewRequest("POST", "/api/v1/setup/exchangeToken", bytes.NewBuffer(exchangeBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	// Verify token exchange
	assert.Equal(t, http.StatusOK, w2.Code)

	var exchangeResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &exchangeResponse)
	assert.NoError(t, err)
	assert.Contains(t, exchangeResponse, "keyUuid")
	assert.Contains(t, exchangeResponse, "apiKey")
	assert.Contains(t, exchangeResponse, "user")
	assert.Equal(t, "testuser", exchangeResponse["user"])

	apiKey := exchangeResponse["apiKey"].(string)
	assert.Contains(t, apiKey, "sk_")

	// Step 3: Verify API key works for authentication
	// This would test using the API key to access a protected endpoint
	// For now, we'll just verify the key format
	assert.Regexp(t, `^sk_[A-Za-z0-9+/]{43}$`, apiKey)
}
