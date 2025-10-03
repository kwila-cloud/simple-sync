package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthErrorScenariosIntegration(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup handlers
	h := handlers.NewTestHandlers(nil)

	// Register routes
	v1 := router.Group("/api/v1")
	v1.POST("/user/generateToken", h.PostUserGenerateToken)
	v1.POST("/user/resetKey", h.PostUserResetKey)
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	t.Run("InsufficientPermissions", func(t *testing.T) {
		// Try to generate token without proper permissions
		generateRequest := map[string]string{
			"user": "testuser",
		}
		requestBody, _ := json.Marshal(generateRequest)

		req, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=testuser", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "sk_insufficient123456789012345678901234567890")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Expected: 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "Unauthorized", response["error"])
	})

	t.Run("NonExistentUser", func(t *testing.T) {
		// Try to generate token for non-existent user
		generateRequest := map[string]string{
			"user": "nonexistent",
		}
		requestBody, _ := json.Marshal(generateRequest)

		req, _ := http.NewRequest("POST", "/api/v1/user/generateToken?user=nonexistent", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "sk_admin123456789012345678901234567890")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Expected: 401 Unauthorized (not 404 to prevent enumeration)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "Unauthorized", response["error"])
	})

	t.Run("ExpiredSetupToken", func(t *testing.T) {
		// Try to exchange an expired setup token
		exchangeRequest := map[string]interface{}{
			"token":       "EXPIRED-TOKEN",
			"description": "Test Client",
		}
		requestBody, _ := json.Marshal(exchangeRequest)

		req, _ := http.NewRequest("POST", "/api/v1/setup/exchangeToken", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Expected: 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "Unauthorized", response["error"])
	})

	t.Run("InvalidTokenFormat", func(t *testing.T) {
		// Try to exchange a malformed token
		exchangeRequest := map[string]interface{}{
			"token":       "INVALID-FORMAT",
			"description": "Test Client",
		}
		requestBody, _ := json.Marshal(exchangeRequest)

		req, _ := http.NewRequest("POST", "/api/v1/setup/exchangeToken", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Expected: 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "Unauthorized", response["error"])
	})
}
