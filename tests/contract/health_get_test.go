package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret")

	// Register routes
	router.GET("/health", h.GetHealth)

	// Test GET /health endpoint
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expected: 200 OK with health response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response structure
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "uptime")

	// Validate status is "healthy"
	assert.Equal(t, "healthy", response["status"])

	// Validate timestamp format (should be ISO 8601)
	timestamp, ok := response["timestamp"].(string)
	assert.True(t, ok, "timestamp should be a string")
	assert.NotEmpty(t, timestamp)

	// Validate version is present
	version, ok := response["version"].(string)
	assert.True(t, ok, "version should be a string")
	assert.NotEmpty(t, version)

	// Validate uptime is a number >= 0
	uptime, ok := response["uptime"].(float64)
	assert.True(t, ok, "uptime should be a number")
	assert.GreaterOrEqual(t, uptime, float64(0))
}
