package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !h.aclService.CheckPermission(userIdStr, ".acl", ".acl.addRule") {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions to update ACL",
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
		if err := validateAclRule(rule); err != nil {
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

// checks if an ACL rule has valid data
func validateAclRule(rule models.AclRule) error {
	if !isValidPattern(rule.User) {
		return errors.New("invalid user pattern")
	}
	if !isValidPattern(rule.Item) {
		return errors.New("invalid item pattern")
	}
	if !isValidPattern(rule.Action) {
		return errors.New("invalid action pattern")
	}
	if rule.Type != "allow" && rule.Type != "deny" {
		return errors.New("type must be either 'allow' or 'deny'")
	}
	return nil
}

// Checks if a pattern has valid wildcard usage (at most one at the end)
func isValidPattern(pattern string) bool {
	if pattern == "" {
		return false
	}
	if pattern == "*" {
		return true
	}
	if containsControlChars(pattern) {
		return false
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return !strings.Contains(prefix, "*")
	}
	return !strings.Contains(pattern, "*")
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
