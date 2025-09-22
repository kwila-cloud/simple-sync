package performance

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"simple-sync/src/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpointPerformance(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers()

	// Register health route
	v1 := router.Group("/api/v1")
	v1.GET("/health", h.GetHealth)

	// Test health endpoint performance
	start := time.Now()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert performance (<10ms)
	assert.Less(t, duration, 10*time.Millisecond, "Health endpoint should respond in less than 10ms")
}
