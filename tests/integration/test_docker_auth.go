package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDockerAuthentication(t *testing.T) {
	// This test assumes the service is running via docker-compose
	// In a real CI environment, this would be run against a test container

	// Skip if not running in Docker environment
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Wait for service to be ready
	time.Sleep(5 * time.Second)

	// Test 1: Get authentication token
	token := getAuthToken(t)

	// Test 2: Access protected endpoint with valid token
	testProtectedEndpoint(t, token, http.StatusOK)

	// Test 3: Access protected endpoint without token
	testProtectedEndpointWithoutToken(t)

	// Test 4: Access protected endpoint with invalid token
	testProtectedEndpointWithInvalidToken(t)
}

func getAuthToken(t *testing.T) string {
	authPayload := map[string]string{
		"username": "testuser",
		"password": "testpass123",
	}

	jsonData, err := json.Marshal(authPayload)
	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/auth/token",
		"application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	token, exists := response["token"]
	assert.True(t, exists, "Response should contain token")
	assert.NotEmpty(t, token, "Token should not be empty")

	return token
}

func testProtectedEndpoint(t *testing.T, token string, expectedStatus int) {
	req, err := http.NewRequest("GET", "http://localhost:8080/events", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, expectedStatus, resp.StatusCode)
}

func testProtectedEndpointWithoutToken(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/events")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func testProtectedEndpointWithInvalidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/events", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer invalid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
