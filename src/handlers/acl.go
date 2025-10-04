package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"simple-sync/src/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostACL handles POST /api/v1/acl for submitting ACL events
func (h *Handlers) PostACL(c *gin.Context) {
	var aclRules []models.AclRule

	// Bind JSON request to ACL rules
	if err := c.ShouldBindJSON(&aclRules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate each ACL rule
	for _, rule := range aclRules {
		if err := validateAclRule(&rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Convert ACL rules to regular events with current timestamp
	var events []models.Event
	currentTime := uint64(time.Now().Unix())

	for _, rule := range aclRules {
		payload := map[string]interface{}{
			"user":   rule.User,
			"item":   rule.Item,
			"action": rule.Action,
		}
		payloadJSON, _ := json.Marshal(payload)

		eventAction := ".acl." + rule.Type

		event := models.Event{
			UUID:      uuid.New().String(),
			Timestamp: currentTime,
			User:      rule.User,
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

// checks if an ACL rule has valid data
func validateAclRule(rule *models.AclRule) error {
	if strings.TrimSpace(rule.User) == "" {
		return errors.New("user is required and cannot be empty")
	}
	if strings.TrimSpace(rule.Item) == "" {
		return errors.New("item is required and cannot be empty")
	}
	if strings.TrimSpace(rule.Action) == "" {
		return errors.New("action is required and cannot be empty")
	}
	if rule.Type != "allow" && rule.Type != "deny" {
		return errors.New("type must be either 'allow' or 'deny'")
	}
	// Check for control characters
	if containsControlChars(rule.User) || containsControlChars(rule.Item) || containsControlChars(rule.Action) {
		return errors.New("user, item, and action cannot contain control characters")
	}
	return nil
}

// checks if string contains control characters
func containsControlChars(s string) bool {
	for _, r := range s {
		if r < 32 || r == 127 {
			return true
		}
	}
	return false
}
