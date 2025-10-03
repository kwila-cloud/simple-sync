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
			User:      "user-123",
			Item:      ".acl",
			Action:    ".acl.allow",
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
	_, _, err = h.AuthService().GenerateApiKey(".root", "Test Key")
	assert.NoError(t, err)

	// Create the target user
	user := &models.User{Id: "user-123"}
	err = store.SaveUser(user)
	assert.NoError(t, err)

	// Generate API key for user
	_, userApiKey, err := h.AuthService().GenerateApiKey("user-123", "User Key")
	assert.NoError(t, err)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/acl", h.PostACL) // Will fail until implemented

	// Test invalid ACL data: missing required field
	invalidACLJSON := `[{
		"user": "",
		"item": "item789",
		"action": "read"
	}]`

	req, _ := http.NewRequest("POST", "/api/v1/acl", bytes.NewBufferString(invalidACLJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", userApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert invalid data is rejected with 400 (will fail until validation implemented)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
