package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAclRootUserBypass(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := handlers.NewTestHandlers(nil)

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

	// Post an event with root user (should bypass ACL)
	event := map[string]interface{}{
		"uuid":      "0199c74f-c696-78f8-833a-82f8cf1f1949",
		"timestamp": 1759985518,
		"user":      ".root",
		"item":      "any-item",
		"action":    "any-action",
		"payload":   "{}",
	}
	eventBody, _ := json.Marshal([]map[string]interface{}{event})

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(eventBody))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("X-API-Key", storage.TestingRootApiKey)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should succeed (root bypass)
	assert.Equal(t, http.StatusOK, postW.Code)
}
