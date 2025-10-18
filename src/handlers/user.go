package handlers

import (
	"log"
	"net/http"
	"simple-sync/src/models"

	"github.com/gin-gonic/gin"
)

// PostUserResetKey handles POST /api/v1/user/resetKey
func (h *Handlers) PostUserResetKey(c *gin.Context) {
	// Check if caller has permission (from middleware)
	callerUserId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	callerUserIdStr, ok := callerUserId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		User string `json:"user" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("PostUserResetKey: invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	userId := request.User
	if userId == "" {
		log.Printf("PostUserResetKey: missing user parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter required"})
		return
	}

	if !h.aclService.CheckPermission(callerUserIdStr, userId, ".user.resetKey") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// Invalidate all existing API keys for the user
	err := h.storage.InvalidateUserApiKeys(userId)
	if err != nil {
		log.Printf("PostUserResetKey: failed to invalidate API keys for user %s: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Log the API call as an internal event
	event := models.NewEvent(
		callerUserIdStr,
		".user."+userId,
		".user.resetKey",
		"{}",
	)
	if err := h.storage.SaveEvents([]models.Event{*event}); err != nil {
		log.Printf("Failed to save reset key event for user %s: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API keys invalidated successfully",
	})
}

// PostUserGenerateToken handles POST /api/v1/user/generateToken
func (h *Handlers) PostUserGenerateToken(c *gin.Context) {
	var request struct {
		User string `json:"user" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("PostUserGenerateToken: invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	userId := request.User
	if userId == "" {
		log.Printf("PostUserGenerateToken: missing user parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter required"})
		return
	}

	// Check if caller has permission (from middleware)
	callerUserId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	callerUserIdStr, ok := callerUserId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if !h.aclService.CheckPermission(callerUserIdStr, userId, ".user.generateToken") {
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
		User:    callerUserIdStr,
		Item:    ".user." + userId,
		Action:  ".user.generateToken",
		Payload: "{}",
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

// PostSetupExchangeToken handles POST /api/v1/user/exchangeToken
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
		log.Printf("Failed to exchange setup token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange setup token"})
		return
	}

	// Log the API call as an internal event
	event := models.Event{
		User:    apiKey.UserID,
		Item:    ".user." + apiKey.UserID,
		Action:  ".user.exchangeToken",
		Payload: "{}",
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
