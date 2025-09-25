package handlers

import (
	"errors"
	"net/http"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handlers contains the HTTP handlers for the API
type Handlers struct {
	storage       storage.Storage
	authService   *services.AuthService
	healthHandler *HealthHandler
}

// NewHandlers creates a new handlers instance
func NewHandlers(storage storage.Storage, encryptionKey, version string) *Handlers {
	return &Handlers{
		storage:       storage,
		authService:   services.NewAuthService(encryptionKey, storage),
		healthHandler: NewHealthHandler(version),
	}
}

// NewTestHandlers creates a new handlers instance with test defaults
func NewTestHandlers() *Handlers {
	return NewTestHandlersWithStorage(storage.NewMemoryStorage())
}

// NewTestHandlersWithStorage creates a new handlers instance with test defaults and custom storage
func NewTestHandlersWithStorage(store storage.Storage) *Handlers {
	return NewHandlers(store, "test-encryption-key-32-bytes-123", "test")
}

// AuthService returns the auth service instance
func (h *Handlers) AuthService() *services.AuthService {
	return h.authService
}

// GetHealth handles GET /health
func (h *Handlers) GetHealth(c *gin.Context) {
	h.healthHandler.GetHealth(c)
}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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

	// Get authenticated user from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Basic validation and set user UUID
	for i := range events {
		if events[i].UUID == "" || events[i].ItemUUID == "" || events[i].Action == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// Enhanced timestamp validation
		if err := validateTimestamp(events[i].Timestamp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Override user ID with authenticated user
		events[i].UserUUID = userID.(string)
	}

	// Save events
	if err := h.storage.SaveEvents(events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Return the saved events
	c.JSON(http.StatusOK, events)
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

// PostUserResetKey handles POST /api/v1/user/resetKey
func (h *Handlers) PostUserResetKey(c *gin.Context) {
	userID := c.Query("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter required"})
		return
	}

	// Check if caller has permission (from middleware)
	callerUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// TODO: Implement proper ACL permission check for .user.resetKey
	// For now, allow all authenticated users (temporary until ACL system is implemented)
	// The .root user should always have access according to the specification
	if callerUserID == ".root" {
		// Allow .root user unrestricted access
	} else {
		// TODO: Check ACL rules for .user.resetKey permission on target user
	}

	// Generate setup token
	setupToken, err := h.authService.GenerateSetupToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Log the API call as an internal event
	event := models.Event{
		UUID:      uuid.New().String(),
		Timestamp: uint64(time.Now().Unix()),
		UserUUID:  callerUserID.(string),
		ItemUUID:  ".user." + userID,
		Action:    ".user.resetKey",
		Payload:   "",
	}
	h.storage.SaveEvents([]models.Event{event})

	c.JSON(http.StatusOK, gin.H{
		"token":     setupToken.Token,
		"expiresAt": setupToken.ExpiresAt,
	})
}

// PostUserGenerateToken handles POST /api/v1/user/generateToken
func (h *Handlers) PostUserGenerateToken(c *gin.Context) {
	userID := c.Query("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter required"})
		return
	}

	// Check if caller has permission (from middleware)
	callerUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// TODO: Implement proper ACL permission check for .user.generateToken
	// For now, allow all authenticated users (temporary until ACL system is implemented)
	// The .root user should always have access according to the specification
	if callerUserID == ".root" {
		// Allow .root user unrestricted access
	} else {
		// TODO: Check ACL rules for .user.generateToken permission on target user
	}

	// Generate setup token
	setupToken, err := h.authService.GenerateSetupToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Log the API call as an internal event
	event := models.Event{
		UUID:      uuid.New().String(),
		Timestamp: uint64(time.Now().Unix()),
		UserUUID:  callerUserID.(string),
		ItemUUID:  ".user." + userID,
		Action:    ".user.generateToken",
		Payload:   "",
	}
	h.storage.SaveEvents([]models.Event{event})

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
		UserUUID:  apiKey.UserID,
		ItemUUID:  ".user." + apiKey.UserID,
		Action:    ".user.exchangeToken",
		Payload:   "",
	}
	h.storage.SaveEvents([]models.Event{event})

	c.JSON(http.StatusOK, gin.H{
		"keyUuid":     apiKey.UUID,
		"apiKey":      plainKey,
		"user":        apiKey.UserID,
		"description": apiKey.Description,
	})
}
