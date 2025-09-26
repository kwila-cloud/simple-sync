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

func TestACLRetrieve(t *testing.T) {
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

	// Post an ACL event
	payload, _ := json.Marshal(map[string]interface{}{
		"user":   "someuser",
		"item":   "someitem",
		"action": "read",
	})
	aclEvent := map[string]interface{}{
		"uuid":      "acl-retrieve-123",
		"timestamp": 1640995200,
		"user":      ".root",
		"item":      ".acl",
		"action":    ".acl.allow",
		"payload":   string(payload),
	}
	aclBody, _ := json.Marshal([]map[string]interface{}{aclEvent})

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(aclBody))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", "Bearer "+adminApiKey)
	postW := httptest.NewRecorder()
	router.ServeHTTP(postW, postReq)
	assert.Equal(t, http.StatusOK, postW.Code)

	// Now, retrieve ACL events
	getReq, _ := http.NewRequest("GET", "/api/v1/events?itemUuid=.acl", nil)
	getReq.Header.Set("Authorization", "Bearer "+adminApiKey)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var responseEvents []map[string]interface{}
	err = json.Unmarshal(getW.Body.Bytes(), &responseEvents)
	assert.NoError(t, err)

	// Should include the ACL event
	var aclEventFound map[string]interface{}
	for _, event := range responseEvents {
		if event["uuid"] == "acl-retrieve-123" {
			aclEventFound = event
			break
		}
	}
	assert.NotNil(t, aclEventFound)
	assert.Equal(t, ".acl", aclEventFound["item"])
	assert.Equal(t, ".acl.allow", aclEventFound["action"])
}
