package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
)

// Handlers contains the HTTP handlers for events
type Handlers struct {
	storage       storage.Storage
	authService   *services.AuthService
	healthHandler *HealthHandler
}

// NewHandlers creates a new handlers instance
func NewHandlers(storage storage.Storage, jwtSecret string) *Handlers {
	return &Handlers{
		storage:       storage,
		authService:   services.NewAuthService(jwtSecret, storage),
		healthHandler: NewHealthHandler("dev"), // Default version, will be overridden
	}
}

// AuthService returns the auth service instance
func (h *Handlers) AuthService() *services.AuthService {
	return h.authService
}

// GetHealth handles GET /health
func (h *Handlers) GetHealth(c *gin.Context) {
	h.healthHandler.GetHealth(c)
}

// SetVersion sets the version for the health handler
func (h *Handlers) SetVersion(version string) {
	h.healthHandler.version = version
}

// GetEvents handles GET /events
func (h *Handlers) GetEvents(c *gin.Context) {
	// Check authenticated user
	_, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Determine fromTimestamp filter (nil means no filter)
	var fromTimestamp *uint64

	fromTimestampStr := c.Query("fromTimestamp")
	if fromTimestampStr != "" {
		parsedTimestamp, err := strconv.ParseUint(fromTimestampStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp parameter"})
			return
		}

		// Validate timestamp bounds
		if err := validateTimestamp(parsedTimestamp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp parameter"})
			return
		}

		fromTimestamp = &parsedTimestamp
	}

	// Load events (filtered or all)
	events, err := h.storage.LoadEvents(fromTimestamp)
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
	userUUID, exists := c.Get("user_uuid")
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

		// Override user UUID with authenticated user
		events[i].UserUUID = userUUID.(string)
	}

	// Save events
	if err := h.storage.SaveEvents(events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Return all events (including newly added)
	allEvents, err := h.storage.LoadEvents(nil)
	if err != nil {
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

// PostAuthToken handles POST /auth/token
func (h *Handlers) PostAuthToken(c *gin.Context) {
	var authRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Authenticate user
	user, err := h.authService.Authenticate(authRequest.Username, authRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
