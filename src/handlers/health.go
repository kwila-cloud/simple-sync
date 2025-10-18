package handlers

import (
	"net/http"
	"time"

	"simple-sync/src/models"

	"github.com/gin-gonic/gin"
)

// GetHealth handles GET /health
func (h Handlers) GetHealth(c *gin.Context) {
	uptime := int64(time.Since(h.startTime).Seconds())

	healthResponse := models.NewHealthCheckResponse("healthy", h.version, uptime)

	c.JSON(http.StatusOK, healthResponse)
}
