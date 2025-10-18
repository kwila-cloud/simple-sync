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
func NewHandlers(storage storage.Storage, version string) (*Handlers, error) {
	authService := services.NewAuthService(storage)
	aclService, err := services.NewAclService(storage)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		storage:     storage,
		authService: authService,
		aclService:  aclService,
		startTime:   time.Now(),
		version:     version,
	}, nil
}

// NewTestHandlers creates a new handlers instance with test defaults
// Panics on error for test convenience
func NewTestHandlers(aclRules []models.AclRule) *Handlers {
	return NewTestHandlersWithStorage(storage.NewTestStorage(aclRules))
}

// NewTestHandlersWithStorage creates a new handlers instance with test defaults and custom storage
// Panics on error for test convenience
func NewTestHandlersWithStorage(store storage.Storage) *Handlers {
	h, err := NewHandlers(store, "test")
	if err != nil {
		panic("Failed to create test handlers: " + err.Error())
	}
	return h
}

// AuthService returns the auth service instance
func (h *Handlers) AuthService() *services.AuthService {
	return h.authService
}

// AclService returns the ACL service instance
func (h *Handlers) AclService() *services.AclService {
	return h.aclService
}
