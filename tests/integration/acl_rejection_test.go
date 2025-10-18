package integration

import (
	"bytes"
	"fmt"
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

func TestAclRejectionViaEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow user to update the ACL
	aclRules := []models.AclRule{
		{
			User:   storage.TestingUserId,
			Item:   ".acl",
			Action: ".acl.*",
			Type:   "allow",
		},
	}

	// Setup handlers with memory storage
	h := handlers.NewTestHandlersOrDie(aclRules)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Attempt to submit ACL event via /events (should be rejected)
	aclEventJSON := fmt.Sprintf(`[{
		"uuid": "acl-test-123",
		"timestamp": 1640995200,
		"user": "%s",
		"item": ".acl",
		"action": ".acl.addRule",
		"payload": "{\"user\":\"user-456\",\"item\":\"item789\",\"action\":\"read\",\"type\":\"allow\"}"
	}]`, storage.TestingUserId)

	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(aclEventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert ACL event is rejected
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
