package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostEventsProtected(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

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

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

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

	// Expected: 403 Forbidden - deny by default
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Insufficient permissions", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}

func TestPostEventsWithInvalidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

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

func TestPostEventsAclValidationFailure(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

	// First, give the user permission to perform ACL operations
	aclEventJSON := `[{
  		"uuid": "acl-perm-uuid",
  		"timestamp": 1640995200,
  		"user": "user-123",
  		"item": ".acl",
  		"action": ".acl.allow",
  		"payload": "{\"user\":\"user-123\",\"item\":\".acl\",\"action\":\".acl.invalid\",\"type\":\"allow\"}"
  	}]`

	aclReq, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(aclEventJSON))
	aclReq.Header.Set("Content-Type", "application/json")
	aclReq.Header.Set("Authorization", "Bearer "+plainKey)
	aclW := httptest.NewRecorder()
	router.ServeHTTP(aclW, aclReq)
	// This should succeed since root user can do anything

	// Test data - ACL event with invalid action
	eventJSON := `[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "user-123",
  		"item": ".acl",
  		"action": ".acl.invalid",
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
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Cannot modify ACL rules", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}

func TestPostEventsMissingRequiredFields(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

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
	err = json.Unmarshal(w.Body.Bytes(), &response)
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
	h := handlers.NewTestHandlers()

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

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
	err = json.Unmarshal(w.Body.Bytes(), &response)
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
	h := handlers.NewTestHandlers()

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

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
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Cannot submit events for other users", response["error"])
	assert.Contains(t, response, "eventUuid")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response["eventUuid"])
}
