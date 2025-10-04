package handlers

import (
	"simple-sync/src/models"
	"simple-sync/src/services"
	"simple-sync/src/storage"
	"time"
)

// Handlers contains the HTTP handlers for the API
type Handlers struct {
	storage     storage.Storage
	authService *services.AuthService
	aclService  *services.AclService
	startTime   time.Time
	version     string
}

// NewHandlers creates a new handlers instance
func NewHandlers(storage storage.Storage, version string) *Handlers {
	return &Handlers{
		storage:     storage,
		authService: services.NewAuthService(storage),
		aclService:  services.NewAclService(storage),
		startTime:   time.Now(),
		version:     version,
	}
}

// NewTestHandlers creates a new handlers instance with test defaults
func NewTestHandlers(aclRules []models.AclRule) *Handlers {
	return NewTestHandlersWithStorage(storage.NewTestStorage(aclRules))
}

// NewTestHandlersWithStorage creates a new handlers instance with test defaults and custom storage
func NewTestHandlersWithStorage(store storage.Storage) *Handlers {
	return NewHandlers(store, "test")
}

// AuthService returns the auth service instance
func (h *Handlers) AuthService() *services.AuthService {
	return h.authService
}

// AclService returns the ACL service instance
func (h *Handlers) AclService() *services.AclService {
	return h.aclService
}
