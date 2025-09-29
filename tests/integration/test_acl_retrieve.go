package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-sync/src/handlers"
	"simple-sync/src/middleware"
	"simple-sync/src/models"
	"simple-sync/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestACLRetrieve(t *testing.T) {
	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup ACL rules
	aclRules := []models.AclRule{
		{
			User:      "someuser",
			Item:      "someitem",
			Action:    "read",
			Type:      "allow",
			Timestamp: 1640995200,
		},
	}

	// Setup handlers with memory storage
	store := storage.NewMemoryStorage(aclRules)
	h := handlers.NewTestHandlersWithStorage(store)

	// Create root user
	rootUser := &models.User{Id: ".root"}
	err := store.SaveUser(rootUser)
	assert.NoError(t, err)

	// Create API key for root
	_, adminApiKey, err := h.AuthService().GenerateApiKey(".root", "Test Key")
	assert.NoError(t, err)

	// Register routes
	v1 := router.Group("/api/v1")

	// Auth routes with middleware
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware(h.AuthService()))
	auth.POST("/user/generateToken", h.PostUserGenerateToken)
	auth.GET("/events", h.GetEvents)
	auth.POST("/events", h.PostEvents)

	// Setup routes
	v1.POST("/setup/exchangeToken", h.PostSetupExchangeToken)

	// Retrieve ACL events
	getReq, _ := http.NewRequest("GET", "/api/v1/events?itemUuid=.acl", nil)
	getReq.Header.Set("Authorization", "Bearer "+adminApiKey)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var responseEvents []map[string]interface{}
	err = json.Unmarshal(getW.Body.Bytes(), &responseEvents)
	assert.NoError(t, err)

	// Should include the ACL event
	var aclEventFound map[string]interface{}
	for _, event := range responseEvents {
		if event["item"] == ".acl" {
			aclEventFound = event
			break
		}
	}
	assert.NotNil(t, aclEventFound)
	assert.Equal(t, ".acl", aclEventFound["item"])
	assert.Equal(t, ".acl.allow", aclEventFound["action"])
}
