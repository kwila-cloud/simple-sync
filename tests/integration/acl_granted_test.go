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

func TestAclPermissionGranted(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow test user to write on allowed-item
	aclRules := []models.AclRule{
		{
			User:   storage.TestingUserId,
			Item:   "allowed-item",
			Action: "write",
			Type:   "allow",
		},
	}

	h := handlers.NewTestHandlers(aclRules)

	// Register routes
	v1 := router.Group("/api/v1")

	// Auth routes with middleware
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Setup routes
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Post an event with permission
	event := map[string]interface{}{
		"uuid":      "0199c74f-c696-78f8-833a-82f8cf1f1949",
		"timestamp": 1759985518,
		"user":      storage.TestingUserId,
		"item":      "allowed-item",
		"action":    "write",
		"payload":   "{}",
	}
	eventBody, _ := json.Marshal([]map[string]interface{}{event})

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(eventBody))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("X-API-Key", storage.TestingApiKey)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should succeed due to ACL allow rule
	assert.Equal(t, http.StatusOK, postW.Code)
}
