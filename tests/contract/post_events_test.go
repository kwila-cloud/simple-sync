package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules to allow the test user to create events
	aclRules := []models.AclRule{
		{
			User:   storage.TestingUserId,
			Item:   "item456",
			Action: "create",
			Type:   "allow",
		},
	}

	// Setup handlers
	h := handlers.NewTestHandlers(aclRules)

	// Register routes with auth middleware
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)

	// Sample event data
	eventJSON := fmt.Sprintf(`[{
  		"uuid": "123e4567-e89b-12d3-a456-426614174000",
  		"timestamp": 1640995200,
  		"user": "%s",
  		"item": "item456",
  		"action": "create",
  		"payload": "{}"
   	}]`, storage.TestingUserId)

	// Create test request
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response - should succeed
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, len(response) >= 1)

	// Find the posted event
	var postedEvent map[string]interface{}
	for _, e := range response {
		if e["uuid"] == "123e4567-e89b-12d3-a456-426614174000" {
			postedEvent = e
			break
		}
	}
	assert.NotNil(t, postedEvent)

	// Check the returned event matches what was posted
	assert.Equal(t, float64(1640995200), postedEvent["timestamp"])
	assert.Equal(t, storage.TestingUserId, postedEvent["user"])
	assert.Equal(t, "item456", postedEvent["item"])
	assert.Equal(t, "create", postedEvent["action"])
	assert.Equal(t, "{}", postedEvent["payload"])
}

func TestConcurrentPostEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes with auth
	v1 := router.Group("/api/v1")
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/events", h.PostEvents)
	auth.GET("/events", h.GetEvents)

	var wg sync.WaitGroup
	numGoroutines := 10
	eventsPerGoroutine := 5
	expectedUUIDs := make(map[string]bool)
	var uuidMutex sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				uuid := fmt.Sprintf("%d-%d", id, j)
				event := fmt.Sprintf(`[{"uuid":"%s","timestamp":%d,"user":"%s","item":"i","action":"a","payload":"p"}]`, uuid, id*100+j+1, storage.TestingUserId)
				req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(event))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-API-Key", storage.TestingApiKey)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusForbidden, w.Code)
				uuidMutex.Lock()
				expectedUUIDs[uuid] = true
				uuidMutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Check total events - should be 0 since all posts were denied
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	req.Header.Set("X-API-Key", storage.TestingApiKey)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var events []models.Event
	err := json.Unmarshal(w.Body.Bytes(), &events)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events))
}
