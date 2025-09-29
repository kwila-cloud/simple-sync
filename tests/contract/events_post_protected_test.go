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

func TestPostEventsProtected(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Test data
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	// Test without Authorization header - should fail with 401
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Authorization header required", response["error"])
}

func TestPostEventsWithValidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow the test user to create item
	aclRules := []models.AclRule{
		{
			User:      "user-123",
			Item:      "item456",
			Action:    "create",
			Type:      "allow",
			Timestamp: 1640995200,
		},
	}

	// Setup handlers
	h := handlers.NewTestHandlers(aclRules)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Use the default API key from memory storage
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Test data - user will be overridden by authenticated user
	eventJSON := `[{
 		"uuid": "123e4567-e89b-12d3-a456-426614174000",
 		"timestamp": 1640995200,
 		"user": "user-123",
 		"item": "item456",
 		"action": "create",
 		"payload": "{}"
 	}]`

	// Test with valid Authorization header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 200 OK with events
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, len(response) >= 1)

	// Find the posted event
	var postedEvent map[string]interface{}
	for _, e := range response {
		if e["uuid"] == "123e4567-e89b-12d3-a456-426614174000" {
			postedEvent = e
			break
		}
	}
	assert.NotNil(t, postedEvent)

	// Check the returned event matches what was posted
	assert.Equal(t, float64(1640995200), postedEvent["timestamp"])
	assert.Equal(t, "user-123", postedEvent["user"])
	assert.Equal(t, "item456", postedEvent["item"])
	assert.Equal(t, "create", postedEvent["action"])
	assert.Equal(t, "{}", postedEvent["payload"])
}

func TestPostEventsWithInvalidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Test data
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "user123",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
  	}]`

	// Test with invalid Authorization header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPostEventsAclPermissionFailure(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Use the default API key from memory storage
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Test data - ACL event (should be denied by default)
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "user-123",
  		"item": ".acl",
  		"action": ".acl.allow",
  		"payload": "{\"user\":\"user2\",\"item\":\"item1\",\"action\":\"read\",\"type\":\"allow\"}"
  	}]`

	// Test with valid Authorization header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 403 Forbidden with eventUuid
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Insufficient permissions", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}

func TestPostEventsMissingRequiredFields(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Use the default API key from memory storage
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Test data - missing action field
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "user-123",
  		"item": "item456",
  		"payload": "{}"
  	}]`

	// Test with valid Authorization header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 400 Bad Request with eventUuid
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Missing required fields", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}

func TestPostEventsInvalidTimestamp(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Use the default API key from memory storage
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Test data - invalid timestamp (zero)
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 0,
  		"user": "user-123",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
  	}]`

	// Test with valid Authorization header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 400 Bad Request with eventUuid
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid timestamp", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}

func TestPostEventsWrongUser(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Use the default API key from memory storage
	plainKey := "sk_ATlUSWpdQVKROfmh47z7q60KjlkQcCaC9ps181Jov8E"

	// Test data - event for different user
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "different-user",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
  	}]`

	// Test with valid Authorization header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 403 Forbidden with eventUuid
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Cannot submit events for other users", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}
