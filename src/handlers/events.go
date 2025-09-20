package handlers

import (
	"net/http"
	"strconv"

	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
)

// Handlers contains the HTTP handlers for events
type Handlers struct {
	storage     *storage.MemoryStorage
	authService *services.AuthService
}

// NewHandlers creates a new handlers instance
func NewHandlers(storage *storage.MemoryStorage, jwtSecret string) *Handlers {
	return &Handlers{
		storage:     storage,
		authService: services.NewAuthService(jwtSecret),
	}
}

// AuthService returns the auth service instance
func (h *Handlers) AuthService() *services.AuthService {
	return h.authService
}

// GetEvents handles GET /events
func (h *Handlers) GetEvents(c *gin.Context) {
	// Check authenticated user
	_, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Check for fromTimestamp query parameter
	fromTimestampStr := c.Query("fromTimestamp")
	if fromTimestampStr != "" {
		fromTimestamp, err := strconv.ParseUint(fromTimestampStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp parameter"})
			return
		}

		// Filter events by timestamp
		allEvents, err := h.storage.LoadEvents()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		filteredEvents := make([]models.Event, 0)
		for _, event := range allEvents {
			if event.Timestamp >= fromTimestamp {
				filteredEvents = append(filteredEvents, event)
			}
		}

		c.JSON(http.StatusOK, filteredEvents)
		return
	}

	// Return all events
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
		if events[i].Timestamp == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp"})
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
	allEvents, err := h.storage.LoadEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, allEvents)
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
