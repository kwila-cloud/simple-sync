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

func TestAclRejectionViaEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow the default user to create events
	aclRules := []models.AclRule{
		{
			User:   "user-123",
			Item:   "item456",
			Action: "create",
			Type:   "allow",
		},
	}

	// Setup handlers with memory storage
	store := storage.NewMemoryStorage(aclRules)
	h := handlers.NewTestHandlersWithStorage(store)

	// Create the default user
	user := &models.User{Id: "user-123"}
	err := store.SaveUser(user)
	assert.NoError(t, err)

	// Use default API key for user-123
	userApiKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Attempt to submit ACL event via /events (should be rejected)
	aclEventJSON := `[{
		"uuid": "acl-test-123",
		"timestamp": 1640995200,
		"user": "user-123",
		"item": ".acl",
		"action": ".acl.allow",
		"payload": "{\"user\":\"user-456\",\"item\":\"item789\",\"action\":\"read\"}"
	}]`

	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(aclEventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", userApiKey)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert ACL event is rejected (will fail until rejection logic implemented)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
