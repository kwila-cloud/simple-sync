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

func TestACLSetup(t *testing.T) {
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
	setupReq.Header.Set("Authorization", "Bearer "+adminApiKey)
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
	_ = exchangeResponse["apiKey"].(string) // API key generated but not used in this test

	// Now, set ACL rule using root API key
	payload, _ := json.Marshal(map[string]interface{}{
		"user":   "testuser",
		"item":   "testitem",
		"action": "write",
	})
	aclEvent := map[string]interface{}{
		"uuid":      "acl-123",
		"timestamp": 1640995200,
		"user":      ".root",
		"item":      ".acl",
		"action":    ".acl.allow",
		"payload":   string(payload),
	}
	aclBody, _ := json.Marshal([]map[string]interface{}{aclEvent})

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(aclBody))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", "Bearer "+adminApiKey) // Use root key
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should succeed
	assert.Equal(t, http.StatusOK, postW.Code)

	// Verify the ACL event was stored
	getReq, _ := http.NewRequest("GET", "/api/v1/events?itemUuid=.acl", nil)
	getReq.Header.Set("Authorization", "Bearer "+adminApiKey)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var responseEvents []map[string]interface{}
	err = json.Unmarshal(getW.Body.Bytes(), &responseEvents)
	assert.NoError(t, err)

	// Find the ACL event
	var aclEventFound map[string]interface{}
	for _, event := range responseEvents {
		if event["uuid"] == "acl-123" {
			aclEventFound = event
			break
		}
	}
	assert.NotNil(t, aclEventFound)
	assert.Equal(t, ".acl", aclEventFound["item"])
	assert.Equal(t, ".acl.allow", aclEventFound["action"])
}
