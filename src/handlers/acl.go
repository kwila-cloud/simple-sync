package handlers

import (
	"encoding/json"
	"net/http"

	"simple-sync/src/models"

	"github.com/gin-gonic/gin"
)

// PostAcl handles POST /api/v1/acl for submitting ACL rules
func (h *Handlers) PostAcl(c *gin.Context) {
	// Get authenticated user from context
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if !h.aclService.CheckPermission(userIdStr, ".acl", ".acl.addRule") {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		return
	}

	var aclRules []models.AclRule

	// Bind JSON request to ACL rules
	if err := c.ShouldBindJSON(&aclRules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate that at least one ACL rule is provided
	if len(aclRules) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one ACL rule required"})
		return
	}

	// Validate each ACL rule
	for _, rule := range aclRules {
		if err := rule.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Convert ACL rules to regular events with current timestamp
	var events []models.Event

	for _, rule := range aclRules {
		ruleJson, _ := json.Marshal(rule)

		events = append(events, *models.NewEvent(
			userId.(string),
			".acl",
			".acl.addRule",
			string(ruleJson),
		))
	}

	// Store the events
	if err := h.storage.SaveEvents(events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ACL events submitted"})
}
