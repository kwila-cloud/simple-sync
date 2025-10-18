package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInvalidCredentialsHandling(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlersOrDie(nil)

	// Register routes
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.GET("/events", h.GetEvents)

	// Test: Access with invalid API key
	getReq, _ := http.NewRequest("GET", "/api/v1/events", nil)
	getReq.Header.Set("X-API-Key", "invalid-api-key")
	getW := httptest.NewRecorder()

	router.ServeHTTP(getW, getReq)

	// Expected: 401
	assert.Equal(t, http.StatusUnauthorized, getW.Code)
}
