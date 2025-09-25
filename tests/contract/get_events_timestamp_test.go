package contract

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEventsWithTimestamp(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

	// Create test request with timestamp query
	req, _ := http.NewRequest("GET", "/api/v1/events?fromTimestamp=1640995200", nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Should return events from timestamp onwards (initially empty)
	expected := "[]"
	assert.JSONEq(t, expected, w.Body.String())
}

func TestGetEventsWithTimestampFiltering(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Generate setup token and exchange for API key
	setupToken, err := h.AuthService().GenerateSetupToken("user-123")
	assert.NoError(t, err)
	_, plainKey, err := h.AuthService().ExchangeSetupToken(setupToken.Token, "test")
	assert.NoError(t, err)

	// Post some events with different timestamps
	eventJSON := `[{"uuid":"1","timestamp":100,"userUuid":"u1","itemUuid":"i1","action":"a","payload":"p"}, {"uuid":"2","timestamp":200,"userUuid":"u2","itemUuid":"i2","action":"b","payload":"q"}, {"uuid":"3","timestamp":300,"userUuid":"u3","itemUuid":"i3","action":"c","payload":"r"}]`

	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Now GET with fromTimestamp=150
	req2, _ := http.NewRequest("GET", "/api/v1/events?fromTimestamp=150", nil)
	req2.Header.Set("Authorization", "Bearer "+plainKey)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Should return events with timestamp >= 150, i.e., 200 and 300
	// Note: userUuid is now overridden with authenticated user's UUID
	expected := `[{"uuid":"2","timestamp":200,"userUuid":"user-123","itemUuid":"i2","action":"b","payload":"q"}, {"uuid":"3","timestamp":300,"userUuid":"user-123","itemUuid":"i3","action":"c","payload":"r"}]`
	assert.JSONEq(t, expected, w2.Body.String())
}
