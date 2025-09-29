package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

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

	// Create test request
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Parse response body as JSON array
	// Initially should be empty array
	expected := "[]"
	assert.JSONEq(t, expected, w.Body.String())
}
