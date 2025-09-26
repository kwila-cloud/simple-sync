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

func TestSuccessfulAuthenticationFlow(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers with memory storage
	store := storage.NewMemoryStorage()
	h := handlers.NewTestHandlersWithStorage(store)

	// Create root user
	rootUser := &models.User{Id: ".root"}
	err := store.SaveUser(rootUser)
	assert.NoError(t, err)

	// Create API key for root
	_, adminApiKey, err := h.AuthService().GenerateApiKey(".root", "Test Key")
	assert.NoError(t, err)

	// Create the target user
	user := &models.User{Id: "user-123"}
	err = store.SaveUser(user)
	assert.NoError(t, err)

	// Set up ACL rule to allow user-123 to perform "create" actions on "item456"
	aclEvent := &models.Event{
		UUID:      "acl-rule-uuid",
		Timestamp: 1640995200,
		User:      ".root",
		Item:      ".acl",
		Action:    ".acl.allow",
		Payload:   `{"user":"user-123","item":"item456","action":"create","type":"allow"}`,
	}
	err = store.SaveEvents([]models.Event{*aclEvent})
	assert.NoError(t, err)

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
	setupReq, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=user-123", nil)
	setupReq.Header.Set("Authorization", "Bearer "+adminApiKey)
	setupW := httptest.NewRecorder()

	router.ServeHTTP(setupW, setupReq)

	assert.Equal(t, http.StatusOK, setupW.Code)

	var setupResponse map[string]string
	err = json.Unmarshal(setupW.Body.Bytes(), &setupResponse)
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
  		"user": "user-123",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
  	}]`

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", "Bearer "+apiKey)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should fail with 403 due to deny-by-default ACL
	assert.Equal(t, http.StatusForbidden, postW.Code)

	var response map[string]interface{}
	err = json.Unmarshal(postW.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Insufficient permissions", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}
