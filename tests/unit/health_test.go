package unit

import (
	"simple-sync/src/models"
	"testing"
)

func TestNewHealthCheckResponse(t *testing.T) {
	status := "healthy"
	version := "1.0.0"
	uptime := int64(123)

	response := models.NewHealthCheckResponse(status, version, uptime)

	if response.Status != status {
		t.Errorf("Expected status %s, got %s", status, response.Status)
	}

	if response.Version != version {
		t.Errorf("Expected version %s, got %s", version, response.Version)
	}

	if response.Uptime != uptime {
		t.Errorf("Expected uptime %d, got %d", uptime, response.Uptime)
	}
}
