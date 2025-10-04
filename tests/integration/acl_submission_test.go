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

func TestACLSubmission(t *testing.T) {
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
	_, adminApiKey, err := h.AuthService().GenerateApiKey(".root", "Test Key")
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
	auth.GET("/events", h.GetEvents)

	// Submit ACL rule
	aclJSON := `[{
		"user": "user-456",
		"item": "item789",
		"action": "read",
		"type": "allow"
	}]`

	req, _ := http.NewRequest("POST", "/api/v1/acl", bytes.NewBufferString(aclJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", userApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert ACL submission succeeds (will fail until implemented)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify ACL event was stored by querying events
	getReq, _ := http.NewRequest("GET", "/api/v1/events?itemUuid=.acl", nil)
	getReq.Header.Set("X-API-Key", adminApiKey)

	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	assert.Equal(t, http.StatusOK, getW.Code)

	var events []models.Event
	err = json.Unmarshal(getW.Body.Bytes(), &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, ".acl", events[0].Item)
	assert.Equal(t, ".acl.allow", events[0].Action)
}
