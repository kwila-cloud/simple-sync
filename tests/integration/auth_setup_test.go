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

func TestUserSetupFlowIntegration(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers with memory storage
	store := storage.NewMemoryStorage()
	h := handlers.NewTestHandlersWithStorage(store)

	// Create root user and API key for authentication
	rootUser := &models.User{Id: ".root"}
	err := store.SaveUser(rootUser)
	assert.NoError(t, err)

	_, adminApiKey, err := h.AuthService().GenerateApiKey(".root", "Admin Key")
	assert.NoError(t, err)

	// Create the target user for token generation
	testUser := &models.User{Id: "testuser"}
	err = store.SaveUser(testUser)
	assert.NoError(t, err)

	// Register routes
	v1 := router.Group("/api/v1")

	// Auth routes with middleware
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/user/generateToken", h.PostUserGenerateToken)

	// Setup routes (no middleware)
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Step 1: Generate setup token for user
	generateRequest := map[string]string{
		"user": "testuser",
	}
	requestBody, _ := json.Marshal(generateRequest)

	req1, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=testuser", bytes.NewBuffer(requestBody))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+adminApiKey)
	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, req1)

	// Verify token generation
	assert.Equal(t, http.StatusOK, w1.Code)

	var generateResponse map[string]string
	err = json.Unmarshal(w1.Body.Bytes(), &generateResponse)
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
