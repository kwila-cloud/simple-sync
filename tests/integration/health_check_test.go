package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	// This test assumes the service is running via docker-compose
	// In a real CI environment, this would be run against a test container

	// Skip if not running in Docker environment
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if service is available
	if !isServiceAvailable() {
		t.Skip("Service not available at localhost:8080, skipping health check integration test")
	}

	// Wait for service to be ready
	time.Sleep(2 * time.Second)

	// Test health endpoint
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Failed to connect to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Should return 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	// Parse response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Validate response structure
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "uptime")

	// Status should be "healthy"
	assert.Equal(t, "healthy", response["status"])

	// Validate timestamp
	timestamp, ok := response["timestamp"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, timestamp)

	// Validate version
	version, ok := response["version"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, version)

	// Validate uptime
	uptime, ok := response["uptime"].(float64)
	assert.True(t, ok)
	assert.GreaterOrEqual(t, uptime, float64(0))
}
