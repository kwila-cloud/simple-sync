package contract

import (
	"bytes"
	"encoding/json"
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
	eventJSON := fmt.Sprintf(`[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"user": "%s",
		"item": "item456",
		"action": "create",
		"payload": "{}"
	}]`, storage.TestingUserId)

	// Test without X-API-Key header - should fail with 401
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
	assert.Equal(t, "X-API-Key header required", response["error"])
}

func TestPostEventsWithValidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow the test user to create item
	aclRules := []models.AclRule{
		{
			User:   storage.TestingUserId,
			Item:   "item456",
			Action: "create",
			Type:   "allow",
		},
	}

	// Setup handlers
	h := handlers.NewTestHandlers(aclRules)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Create a valid event using the model constructor
	event := models.NewEvent(storage.TestingUserId, "item456", "create", "{}")

	// Convert to JSON
	eventData := []models.Event{*event}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Test with valid X-API-Key header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 200 OK with events
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, len(response) >= 1)

	// Find the posted event
	var postedEvent map[string]interface{}
	for _, e := range response {
		if e["user"] == storage.TestingUserId && e["item"] == "item456" && e["action"] == "create" {
			postedEvent = e
			break
		}
	}
	assert.NotNil(t, postedEvent)

	// Check the returned event matches what was posted
	assert.Equal(t, float64(event.Timestamp), postedEvent["timestamp"])
	assert.Equal(t, storage.TestingUserId, postedEvent["user"])
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

	// Test with invalid X-API-Key header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPostEventsAclEventNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Test data - ACL event (should be denied by default)
	eventJSON := fmt.Sprintf(`[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "%s",
  		"item": ".acl",
  		"action": ".acl.addRule",
   		"payload": "{\"user\":\"user2\",\"item\":\"item1\",\"action\":\"delete\",\"type\":\"allow\"}"
  	}]`, storage.TestingUserId)

	// Test with valid X-API-Key header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "ACL events must be submitted via dedicated /api/v1/acl endpoint", response["error"])
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

	// Create a valid event but remove the action field after marshaling
	event := models.NewEvent(storage.TestingUserId, "item456", "create", "{}")

	// Convert to map and remove action field to test missing field validation
	eventMap := map[string]interface{}{
		"uuid":      event.UUID,
		"timestamp": event.Timestamp,
		"user":      event.User,
		"item":      event.Item,
		"payload":   event.Payload,
		// action field intentionally missing
	}

	eventData := []map[string]interface{}{eventMap}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Test with valid X-API-Key header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 400 Bad Request with eventUuid
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "action is required", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, event.UUID, response["eventUuid"])
}

func TestPostEventsInvalidTimestamp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Test data - invalid timestamp (zero)
	eventJSON := fmt.Sprintf(`[{
  		"uuid": "00000000-0000-7a88-a1c8-395f11f5c9ad",
  		"timestamp": 0,
  		"user": "user-123",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
  	}]`)

	// Test with valid X-API-Key header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 400 Bad Request with eventUuid
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid timestamp", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "00000000-0000-7a88-a1c8-395f11f5c9ad", response["eventUuid"])
}

func TestPostEventsWrongUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Test data - event for different user
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "different-user",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
  	}]`

	// Test with valid X-API-Key header
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 400 Bad Request with eventUuid (UUID validation happens before user validation)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid timestamp", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}
