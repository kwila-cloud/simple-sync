package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulAuthenticationFlow(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store, "test-secret")

	// Register routes
	router.POST("/auth/token", h.PostAuthToken)
	router.GET("/events", h.GetEvents)
	router.POST("/events", h.PostEvents)

	// Step 1: Authenticate and get token
	authRequest := map[string]string{
		"username": "testuser",
		"password": "testpass123",
	}
	authBody, _ := json.Marshal(authRequest)

	authReq, _ := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(authBody))
	authReq.Header.Set("Content-Type", "application/json")
	authW := httptest.NewRecorder()

	router.ServeHTTP(authW, authReq)

	// Should get token (will fail until implemented)
	assert.Equal(t, http.StatusOK, authW.Code)

	var authResponse map[string]string
	err := json.Unmarshal(authW.Body.Bytes(), &authResponse)
	assert.NoError(t, err)
	token := authResponse["token"]
	assert.NotEmpty(t, token)

	// Step 2: Use token to access protected GET /events
	getReq, _ := http.NewRequest("GET", "/events", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getW := httptest.NewRecorder()

	router.ServeHTTP(getW, getReq)

	// Should succeed (will fail until middleware implemented)
	assert.Equal(t, http.StatusOK, getW.Code)

	// Step 3: Use token to POST events
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	postReq, _ := http.NewRequest("POST", "/events", bytes.NewBufferString(eventJSON))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", "Bearer "+token)
	postW := httptest.NewRecorder()

	router.ServeHTTP(postW, postReq)

	// Should succeed
	assert.Equal(t, http.StatusOK, postW.Code)
	assert.JSONEq(t, eventJSON, postW.Body.String())
}
