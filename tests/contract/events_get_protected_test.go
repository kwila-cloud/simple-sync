package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEventsProtected(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Test without X-API-Key header - should fail with 401
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
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

func TestGetEventsWithValidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

	// Test with valid X-API-Key header
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	req.Header.Set("X-API-Key", plainKey)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 200 OK with events
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Should return empty array initially
	assert.JSONEq(t, "[]", w.Body.String())
}

func TestGetEventsWithInvalidToken(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Test with invalid X-API-Key header
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	req.Header.Set("X-API-Key", "invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
