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

func TestAuthenticationFlow(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow the test user to create events
	aclRules := []models.AclRule{
		{
			User:   storage.TestingUserId,
			Item:   "item456",
			Action: "create",
			Type:   "allow",
		},
	}

	// Setup handlers with memory storage
	h := handlers.NewTestHandlersOrDie(aclRules)

	// Register routes
	v1 := router.Group("/api/v1")

	// Auth routes with middleware
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/user/generateToken", h.PostUserGenerateToken)
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Setup routes
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Step 1: Generate setup token
	generateRequest := map[string]interface{}{
		"user": storage.TestingRootApiKey,
	}
	requestBody, _ := json.Marshal(generateRequest)
	setupReq, _ := http.NewRequest("POST", "/api/v1/user/generateToken", bytes.NewBuffer(requestBody))
	setupReq.Header.Set("Content-Type", "application/json")
	setupReq.Header.Set("X-API-Key", storage.TestingRootApiKey)
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

	exchangeReq, _ := http.NewRequest("POST", "/api/v1/user/exchangeToken", bytes.NewBuffer(exchangeBody))
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
	getReq.Header.Set("X-API-Key", apiKey)
	getW := httptest.NewRecorder()

	router.ServeHTTP(getW, getReq)

	assert.Equal(t, http.StatusOK, getW.Code)

	// Step 4: Use API key to POST events
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "` + storage.TestingRootApiKey + `",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
  	}]`

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("X-API-Key", apiKey)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should succeed
	assert.Equal(t, http.StatusOK, postW.Code)

	var response []map[string]interface{}
	err = json.Unmarshal(postW.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, len(response) >= 1)

	// Find the posted event
	var postedEvent map[string]interface{}
	for _, e := range response {
		if e["uuid"] == "123e4567-e89b-12d3-a456-426614174000" {
			postedEvent = e
			break
		}
	}
	assert.NotNil(t, postedEvent)

	// Check the returned event matches what was posted
	assert.Equal(t, float64(1640995200), postedEvent["timestamp"])
	assert.Equal(t, storage.TestingUserId, postedEvent["user"])
	assert.Equal(t, "item456", postedEvent["item"])
	assert.Equal(t, "create", postedEvent["action"])
	assert.Equal(t, "{}", postedEvent["payload"])
}
