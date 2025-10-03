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

func TestACLPermissionGranted(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow testuser to write on allowed-item
	aclRules := []models.AclRule{
		{
			User:      "testuser",
			Item:      "allowed-item",
			Action:    "write",
			Type:      "allow",
			Timestamp: 1640995200,
		},
	}

	// Setup handlers with memory storage
	store := storage.NewMemoryStorage(aclRules)
	h := handlers.NewTestHandlersWithStorage(store)

	// Create root user
	rootUser := &models.User{Id: ".root"}
	err := store.SaveUser(rootUser)
	assert.NoError(t, err)

	// Create API key for root
	_, adminApiKey, err := h.AuthService().GenerateApiKey(".root", "Test Key")
	assert.NoError(t, err)

	// Create the target user
	user := &models.User{Id: "testuser"}
	err = store.SaveUser(user)
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

	// Generate API key for testuser
	setupReq, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=testuser", nil)
	setupReq.Header.Set("X-API-Key", adminApiKey)
	setupW := httptest.NewRecorder()
	router.ServeHTTP(setupW, setupReq)
	assert.Equal(t, http.StatusOK, setupW.Code)

	var setupResponse map[string]string
	err = json.Unmarshal(setupW.Body.Bytes(), &setupResponse)
	assert.NoError(t, err)
	setupToken := setupResponse["token"]

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

	// Post an event with permission
	event := map[string]interface{}{
		"uuid":      "event-456",
		"timestamp": 1640995200,
		"user":      "testuser",
		"item":      "allowed-item",
		"action":    "write",
		"payload":   "{}",
	}
	eventBody, _ := json.Marshal([]map[string]interface{}{event})

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(eventBody))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("X-API-Key", apiKey)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should succeed due to ACL allow rule
	assert.Equal(t, http.StatusOK, postW.Code)
}
