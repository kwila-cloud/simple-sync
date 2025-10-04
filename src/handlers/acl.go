package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"simple-sync/src/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostACL handles POST /api/v1/acl for submitting ACL events
func (h *Handlers) PostACL(c *gin.Context) {
	var aclEvents []models.ACLEvent

	// Bind JSON request to ACL events
	if err := c.ShouldBindJSON(&aclEvents); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate each ACL event
	for _, aclEvent := range aclEvents {
		if err := aclEvent.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Convert ACL events to regular events with current timestamp
	var events []models.Event
	currentTime := uint64(time.Now().Unix())

	for _, aclEvent := range aclEvents {
		payload := map[string]interface{}{
			"user":   aclEvent.User,
			"item":   aclEvent.Item,
			"action": aclEvent.Action,
		}
		payloadJSON, _ := json.Marshal(payload)

		eventAction := ".acl." + aclEvent.Type

		event := models.Event{
			UUID:      uuid.New().String(),
			Timestamp: currentTime,
			User:      aclEvent.User,
			Item:      ".acl",
			Action:    eventAction,
			Payload:   string(payloadJSON),
		}
		events = append(events, event)
	}

	// Store the events
	if err := h.storage.SaveEvents(events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ACL events submitted"})
}
