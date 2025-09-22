package models

import "time"

// HealthCheckResponse represents the response structure for health check endpoints
type HealthCheckResponse struct {
	Status    string `json:"status"`    // "healthy" or "unhealthy"
	Timestamp string `json:"timestamp"` // ISO 8601 timestamp
	Version   string `json:"version"`   // Application version
	Uptime    int64  `json:"uptime"`    // Uptime in seconds
}

// NewHealthCheckResponse creates a new health check response
func NewHealthCheckResponse(status, version string, uptime int64) *HealthCheckResponse {
	return &HealthCheckResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
		Uptime:    uptime,
	}
}
