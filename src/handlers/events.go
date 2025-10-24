package handlers

import (
	"log"
	"net/http"

	"simple-sync/src/models"

	"github.com/gin-gonic/gin"
)

// GetEvents handles GET /events
func (h *Handlers) GetEvents(c *gin.Context) {
	// Check authenticated user
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Load all events
	events, err := h.storage.LoadEvents()
	if err != nil {
		log.Printf("GetEvents: failed to load events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// PostEvents handles POST /events
func (h *Handlers) PostEvents(c *gin.Context) {
	// Get authenticated user from context
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Bind JSON array
	var events []models.Event
	if err := c.ShouldBindJSON(&events); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Reject ACL events submitted via /events
	for _, event := range events {
		if event.Item == ".acl" && len(event.Action) > 4 && event.Action[:5] == ".acl." {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ACL events must be submitted via dedicated /api/v1/acl endpoint", "eventUuid": event.UUID})
			return
		}
	}

	// Validate each event using the model validation
	for _, event := range events {
		if err := event.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "eventUuid": event.UUID})
			return
		}

		// Validate that the event user matches the authenticated user
		if event.User != "" && event.User != userId.(string) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot submit events for other users", "eventUuid": event.UUID})
			return
		}
	}

	// ACL permission checks for each event
	for _, event := range events {
		if !h.aclService.CheckPermission(userId.(string), event.Item, event.Action) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions", "eventUuid": event.UUID})
			return
		}
		if event.IsApiOnlyEvent() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot add internal events through this endpoint", "eventUuid": event.UUID})
			return
		}
	}

	// Add events
	if err := h.storage.AddEvents(events); err != nil {
		log.Printf("PostEvents: failed to save events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Return all events (including newly added)
	allEvents, err := h.storage.LoadEvents()
	if err != nil {
		log.Printf("PostEvents: failed to load all events after save: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, allEvents)
}
