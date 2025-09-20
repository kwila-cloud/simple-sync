package handlers

import (
	"net/http"
	"strconv"

	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
)

// Handlers contains the HTTP handlers for events
type Handlers struct {
	storage *storage.MemoryStorage
}

// NewHandlers creates a new handlers instance
func NewHandlers(storage *storage.MemoryStorage) *Handlers {
	return &Handlers{
		storage: storage,
	}
}

// GetEvents handles GET /events
func (h *Handlers) GetEvents(c *gin.Context) {
	// Check for fromTimestamp query parameter
	fromTimestampStr := c.Query("fromTimestamp")
	if fromTimestampStr != "" {
		fromTimestamp, err := strconv.ParseUint(fromTimestampStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fromTimestamp"})
			return
		}

		// Filter events by timestamp
		allEvents, err := h.storage.LoadEvents()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load events"})
			return
		}

		filteredEvents := make([]models.Event, 0)
		for _, event := range allEvents {
			if event.Timestamp >= fromTimestamp {
				filteredEvents = append(filteredEvents, event)
			}
		}

		c.JSON(http.StatusOK, filteredEvents)
		return
	}

	// Return all events
	events, err := h.storage.LoadEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// PostEvents handles POST /events
func (h *Handlers) PostEvents(c *gin.Context) {
	var events []models.Event

	// Bind JSON array
	if err := c.ShouldBindJSON(&events); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Basic validation
	for _, event := range events {
		if event.UUID == "" || event.UserUUID == "" || event.ItemUUID == "" || event.Action == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}
		if event.Timestamp == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp"})
			return
		}
	}

	// Save events
	if err := h.storage.SaveEvents(events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save events"})
		return
	}

	// Return all events (including newly added)
	allEvents, err := h.storage.LoadEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load events"})
		return
	}

	c.JSON(http.StatusOK, allEvents)
}
