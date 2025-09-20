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
	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store)

	// Register routes
	router.POST("/events", h.PostEvents)

	// Sample event data
	eventJSON := `[{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"timestamp": 1640995200,
		"userUuid": "user123",
		"itemUuid": "item456",
		"action": "create",
		"payload": "{}"
	}]`

	// Create test request
	req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Should return the posted events
	assert.JSONEq(t, eventJSON, w.Body.String())
}

func TestConcurrentPostEvents(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup storage and handlers
	store := storage.NewMemoryStorage()
	h := handlers.NewHandlers(store)

	// Register routes
	router.POST("/events", h.PostEvents)
	router.GET("/events", h.GetEvents)

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
				event := fmt.Sprintf(`[{"uuid":"%s","timestamp":%d,"userUuid":"u","itemUuid":"i","action":"a","payload":"p"}]`, uuid, id*100+j+1)
				req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString(event))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Code)
				uuidMutex.Lock()
				expectedUUIDs[uuid] = true
				uuidMutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Check total events
	req, _ := http.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var events []models.Event
	err := json.Unmarshal(w.Body.Bytes(), &events)
	assert.NoError(t, err)
	assert.Equal(t, numGoroutines*eventsPerGoroutine, len(events))

	// Verify all expected UUIDs are present and unique
	actualUUIDs := make(map[string]int)
	for _, event := range events {
		actualUUIDs[event.UUID]++
	}

	// Check that each expected UUID appears exactly once
	for expectedUUID := range expectedUUIDs {
		count, exists := actualUUIDs[expectedUUID]
		assert.True(t, exists, "Expected UUID %s not found in retrieved events", expectedUUID)
		assert.Equal(t, 1, count, "UUID %s appears %d times, expected exactly 1", expectedUUID, count)
	}
}
