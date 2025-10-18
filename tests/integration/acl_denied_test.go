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

func TestACLPermissionDenied(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := handlers.NewTestHandlers(nil)

	// Register routes
	v1 := router.Group("/api/v1")

	// Auth routes with middleware
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Now, try to post an event without permission (deny by default)
	event := map[string]interface{}{
		"uuid":      "event-123",
		"timestamp": 1640995200,
		"user":      storage.TestingUserId,
		"item":      "restricted-item",
		"action":    "write",
		"payload":   "{}",
	}
	eventBody, _ := json.Marshal([]map[string]interface{}{event})

	postReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(eventBody))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("X-API-Key", storage.TestingApiKey)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should be forbidden due to ACL (deny by default)
	assert.Equal(t, http.StatusForbidden, postW.Code)
}
