package handlers

import (
	"errors"
	"log"
	"net/http"
	"simple-sync/src/models"
	"time"

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

	// Basic validation for each event first
	for _, event := range events {
		if event.UUID == "" || event.Item == "" || event.Action == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields", "eventUuid": event.UUID})
			return
		}

		// Enhanced timestamp validation
		if err := validateTimestamp(event.Timestamp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp", "eventUuid": event.UUID})
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
		// For ACL events, additional validation
		if event.IsAclEvent() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify ACL rules through this endpoint", "eventUuid": event.UUID})
			return
		}
	}

	// Save events
	if err := h.storage.SaveEvents(events); err != nil {
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

// validateTimestamp performs enhanced timestamp validation
func validateTimestamp(timestamp uint64) error {
	// Basic zero check
	if timestamp == 0 {
		return errors.New("Invalid timestamp")
	}

	// Maximum timestamp: Allow up to 24 hours in the future for clock skew tolerance
	now := time.Now().Unix()
	maxTimestamp := now + (24 * 60 * 60) // 24 hours from now
	if int64(timestamp) > maxTimestamp {
		return errors.New("Invalid timestamp")
	}

	return nil
}
