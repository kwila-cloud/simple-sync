package handlers

import (
	"log"
	"net/http"
	"simple-sync/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostUserResetKey handles POST /api/v1/user/resetKey
func (h *Handlers) PostUserResetKey(c *gin.Context) {
	userID := c.Query("user")
	if userID == "" {
		log.Printf("PostUserResetKey: missing user parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter required"})
		return
	}

	// Check if caller has permission (from middleware)
	callerUserID, exists := c.Get("user_id")
	if !exists {
		log.Printf("PostUserResetKey: user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// TODO(#5): Implement proper ACL permission check for .user.resetKey
	// For now, allow all authenticated users (temporary until ACL system is implemented)
	// The .root user should always have access according to the specification
	if callerUserID == ".root" {
		// Allow .root user unrestricted access
	} else {
		// TODO(#5): Check ACL rules for .user.resetKey permission on target user
	}

	// Invalidate all existing API keys for the user
	err := h.storage.InvalidateUserAPIKeys(userID)
	if err != nil {
		log.Printf("PostUserResetKey: failed to invalidate API keys for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Log the API call as an internal event
	event := models.Event{
		UUID:      uuid.New().String(),
		Timestamp: uint64(time.Now().Unix()),
		User:      callerUserID.(string),
		Item:      ".user." + userID,
		Action:    ".user.resetKey",
		Payload:   "{}",
	}
	if err := h.storage.SaveEvents([]models.Event{event}); err != nil {
		log.Printf("Failed to save reset key event for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API keys invalidated successfully",
	})
}

// PostUserGenerateToken handles POST /api/v1/user/generateToken
func (h *Handlers) PostUserGenerateToken(c *gin.Context) {
	userId := c.Query("user")
	if userId == "" {
		log.Printf("PostUserGenerateToken: missing user parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter required"})
		return
	}

	// Check if caller has permission (from middleware)
	callerUserId, exists := c.Get("user_id")
	if !exists {
		log.Printf("PostUserGenerateToken: user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.aclService.CheckPermission(callerUserId.(string), userId, ".user.generateToken") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// Generate setup token
	setupToken, err := h.authService.GenerateSetupToken(userId)
	if err != nil {
		log.Printf("PostUserGenerateToken: failed to generate setup token for user %s: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Log the API call as an internal event
	event := models.Event{
		UUID:      uuid.New().String(),
		Timestamp: uint64(time.Now().Unix()),
		User:      callerUserId.(string),
		Item:      ".user." + userId,
		Action:    ".user.generateToken",
		Payload:   "{}",
	}
	if err := h.storage.SaveEvents([]models.Event{event}); err != nil {
		log.Printf("Failed to save generate token event for user %s: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     setupToken.Token,
		"expiresAt": setupToken.ExpiresAt,
	})
}

// PostSetupExchangeToken handles POST /api/v1/setup/exchangeToken
func (h *Handlers) PostSetupExchangeToken(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Exchange setup token for API key
	apiKey, plainKey, err := h.authService.ExchangeSetupToken(request.Token, request.Description)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Log the API call as an internal event
	event := models.Event{
		UUID:      uuid.New().String(),
		Timestamp: uint64(time.Now().Unix()),
		User:      apiKey.UserID,
		Item:      ".user." + apiKey.UserID,
		Action:    ".user.exchangeToken",
		Payload:   "{}",
	}
	if err := h.storage.SaveEvents([]models.Event{event}); err != nil {
		log.Printf("Failed to save exchange token event for user %s: %v", apiKey.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keyUuid":     apiKey.UUID,
		"apiKey":      plainKey,
		"user":        apiKey.UserID,
		"description": apiKey.Description,
	})
}
