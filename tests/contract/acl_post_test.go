package contract

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

func TestPostAcl(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow the test user to submit ACL events
	aclRules := []models.AclRule{
		{
			User:   storage.TestingUserId,
			Item:   ".acl",
			Action: ".acl.addRule",
			Type:   "allow",
		},
	}

	// Setup handlers
	h := handlers.NewTestHandlers(aclRules)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/acl", h.PostAcl)

	// Sample ACL rule data
	aclJSON := `[{
		"user": "user-456",
		"item": "item789",
		"action": "read",
		"type": "allow"
	}]`

	// Create request
	req, _ := http.NewRequest("POST", "/api/v1/acl", bytes.NewBufferString(aclJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ACL events submitted", response["message"])
}

func TestPostAclInsufficientPermissions(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers WITHOUT ACL rules (user has no permission to set ACL rules)
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/acl", h.PostAcl)

	// Sample ACL rule data
	aclJSON := `[{
		"user": "user-456",
		"item": "item789",
		"action": "read",
		"type": "allow"
	}]`

	// Create request
	req, _ := http.NewRequest("POST", "/api/v1/acl", bytes.NewBufferString(aclJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Insufficient permissions to update ACL", response["error"])
}

func TestPostAclInvalidApiKey(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/acl", h.PostAcl)

	// Sample ACL rule data
	aclJSON := `[{
		"user": "user-456",
		"item": "item789",
		"action": "read",
		"type": "allow"
	}]`

	// Create request with invalid API key
	req, _ := http.NewRequest("POST", "/api/v1/acl", bytes.NewBufferString(aclJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "invalid-api-key")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response - should be unauthorized due to invalid API key
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid API key", response["error"])
}
