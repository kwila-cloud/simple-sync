package handlers

import (
	"net/http"
	"time"

	"simple-sync/src/models"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// GetHealth returns the service health status
func (hh *HealthHandler) GetHealth(c *gin.Context) {
	uptime := int64(time.Since(hh.startTime).Seconds())

	healthResponse := models.NewHealthCheckResponse("healthy", hh.version, uptime)

	c.JSON(http.StatusOK, healthResponse)
}
