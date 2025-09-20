package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register routes when handlers are implemented
	// router.GET("/events", handlers.GetEvents)

	// Create test request
	req, _ := http.NewRequest("GET", "/events", nil)
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