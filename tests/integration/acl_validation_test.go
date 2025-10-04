package integration

import (
	"bytes"
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

func TestACLInvalidDataHandling(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow the test user to submit ACL events
	aclRules := []models.AclRule{
		{
			User:   storage.TestingUserId,
			Item:   ".acl",
			Action: ".acl.allow",
			Type:   "allow",
		},
	}

	// Setup handlers with memory storage
	store := storage.NewTestStorage(aclRules)
	h := handlers.NewTestHandlersWithStorage(store)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/acl", h.PostAcl)

	// Test invalid ACL data: missing required field
	invalidACLJSON := `[{
		"user": "",
		"item": "item789",
		"action": "read",
		"type": "allow"
	}]`

	req, _ := http.NewRequest("POST", "/api/v1/acl", bytes.NewBufferString(invalidACLJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
