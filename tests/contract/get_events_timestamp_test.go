package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEventsWithTimestamp(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register routes when handlers are implemented
	// router.GET("/events", handlers.GetEvents)

	// Create test request with timestamp query
	req, _ := http.NewRequest("GET", "/events?fromTimestamp=1640995200", nil)
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