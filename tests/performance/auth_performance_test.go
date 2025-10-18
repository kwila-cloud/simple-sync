package performance

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEventsEndpointPerformance(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()

	start := time.Now()
	router.ServeHTTP(w, req)
	duration := time.Since(start)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Less(t, duration, 200*time.Millisecond, "Get events endpoint should respond in less than 200ms")
}
