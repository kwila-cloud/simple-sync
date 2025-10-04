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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostACL(t *testing.T) {
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

	// Setup handlers
	h := handlers.NewTestHandlers(aclRules)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/acl", h.PostACL) // This will fail until implemented

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
	req.Header.Set("X-API-Key", "test-api-key")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response (will fail until endpoint is implemented)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ACL events submitted", response["message"])
}
