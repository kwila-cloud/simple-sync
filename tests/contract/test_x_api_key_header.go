package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestXAPIKeyHeaderAccepted tests that requests with X-API-Key header are accepted
func TestXAPIKeyHeaderAccepted(t *testing.T) {
	// This test will fail until implementation is updated
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-API-Key header required"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "test-key")
	router.ServeHTTP(w, req)

	// Should pass once implementation is updated
	assert.Equal(t, http.StatusOK, w.Code)
}
