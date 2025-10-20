package storage

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"simple-sync/src/models"
)

const TestingUserId = "user-123"
const TestingApiKey = "sk_fSiYfCABeWUsDjgU3ExViC7/UCkccpxyllbCNJsMGYk"
const TestingRootApiKey = "sk_RYQR7tiqy82dbNEcAdHtO4mbl4YFo9GDF2sr0PbwTlY"

// TestStorage implements in-memory storage for testing,
type TestStorage struct {
	events      []models.Event
	users       map[string]*models.User       // id -> user
	apiKeys     map[string]*models.ApiKey     // uuid -> api key
	setupTokens map[string]*models.SetupToken // token -> setup token
	mutex       sync.RWMutex
}

// NewTestStorage creates a new instance of TestStorage with a default user
func NewTestStorage(aclRules []models.AclRule) *TestStorage {
	storage := &TestStorage{
		events:      make([]models.Event, 0),
		users:       make(map[string]*models.User),
		apiKeys:     make(map[string]*models.ApiKey),
		setupTokens: make(map[string]*models.SetupToken),
	}

	// Add root user
	rootUser, _ := models.NewUser(".root")
	storage.SaveUser(rootUser)

	// Add root API key
	keyHash, _ := bcrypt.GenerateFromPassword([]byte(TestingRootApiKey), bcrypt.DefaultCost)
	now := time.Now()
	apiKey := &models.ApiKey{
		UUID:        "test-root-api-key-uuid",
		UserID:      ".root",
		KeyHash:     string(keyHash),
		Description: "Test Root API Key",
		CreatedAt:   now,
		LastUsedAt:  &now,
	}
	storage.CreateApiKey(apiKey)

	// Add default user
	defaultUser, _ := models.NewUser(TestingUserId)
	storage.SaveUser(defaultUser)

	keyHash, _ = bcrypt.GenerateFromPassword([]byte(TestingApiKey), bcrypt.DefaultCost)
	now = time.Now()
	apiKey = &models.ApiKey{
		UUID:        "test-api-key-uuid",
		UserID:      TestingUserId,
		KeyHash:     string(keyHash),
		Description: "Test API Key",
		CreatedAt:   now,
		LastUsedAt:  &now,
	}
	storage.CreateApiKey(apiKey)

	// Add initial ACL rules as events
	for _, rule := range aclRules {
		ruleJson, _ := json.Marshal(rule)

		storage.events = append(storage.events, *models.NewEvent(
			".root",
			".acl",
			".acl.addRule",
			string(ruleJson),
		))
	}

	return storage
}

// SaveEvents appends new events to the storage
func (m *TestStorage) SaveEvents(events []models.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.events = append(m.events, events...)
	return nil
}

// LoadEvents returns all stored events
func (m *TestStorage) LoadEvents() ([]models.Event, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external modification
	allEvents := make([]models.Event, len(m.events))
	copy(allEvents, m.events)

	return allEvents, nil
}

// SaveUser stores a user by id
func (m *TestStorage) SaveUser(user *models.User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.users[user.Id] = user
	return nil
}

// GetUserById retrieves a user by id
func (m *TestStorage) GetUserById(id string) (*models.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	user, exists := m.users[id]
	if !exists {
		return nil, ErrNotFound
	}
	return user, nil
}

// CreateApiKey stores a new API key
func (m *TestStorage) CreateApiKey(apiKey *models.ApiKey) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.apiKeys[apiKey.UUID] = apiKey
	return nil
}

// GetApiKeyByHash retrieves an API key by its hash
func (m *TestStorage) GetApiKeyByHash(hash string) (*models.ApiKey, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, apiKey := range m.apiKeys {
		if apiKey.KeyHash == hash {
			return apiKey, nil
		}
	}
	return nil, ErrApiKeyNotFound
}

// GetAllApiKeys retrieves all API keys
func (m *TestStorage) GetAllApiKeys() ([]*models.ApiKey, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	keys := make([]*models.ApiKey, 0, len(m.apiKeys))
	for _, k := range m.apiKeys {
		keys = append(keys, k)
	}
	return keys, nil
}

// UpdateApiKey updates an existing API key
func (m *TestStorage) UpdateApiKey(apiKey *models.ApiKey) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.apiKeys[apiKey.UUID] = apiKey
	return nil
}

// CreateSetupToken stores a new setup token
func (m *TestStorage) CreateSetupToken(token *models.SetupToken) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.setupTokens[token.Token] = token
	return nil
}

// GetSetupToken retrieves a setup token by its value
func (m *TestStorage) GetSetupToken(token string) (*models.SetupToken, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	setupToken, exists := m.setupTokens[token]
	if !exists {
		return nil, ErrSetupTokenNotFound
	}
	return setupToken, nil
}

// UpdateSetupToken updates an existing setup token
func (m *TestStorage) UpdateSetupToken(token *models.SetupToken) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.setupTokens[token.Token] = token
	return nil
}

// InvalidateUserSetupTokens marks all setup tokens for a user as used
func (m *TestStorage) InvalidateUserSetupTokens(userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	now := time.Now()
	for _, token := range m.setupTokens {
		if token.UserID == userID {
			token.UsedAt = now
		}
	}
	return nil
}

// InvalidateUserApiKeys removes all API keys for a user
func (m *TestStorage) InvalidateUserApiKeys(userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for uuid, apiKey := range m.apiKeys {
		if apiKey.UserID == userID {
			delete(m.apiKeys, uuid)
		}
	}
	return nil
}

// CreateAclRule stores a new ACL rule
func (m *TestStorage) CreateAclRule(rule *models.AclRule) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Convert ACL rule to event for backward compatibility during transition
	ruleJson, _ := json.Marshal(rule)
	event := models.NewEvent(
		".root",
		".acl",
		".acl.addRule",
		string(ruleJson),
	)
	m.events = append(m.events, *event)

	return nil
}

// GetAclRules retrieves all ACL rules
func (m *TestStorage) GetAclRules() ([]models.AclRule, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var rules []models.AclRule
	for _, event := range m.events {
		if event.IsAclEvent() {
			rule, err := event.ToAclRule()
			if err != nil {
				return nil, fmt.Errorf("malformed ACL rule in event: %w", err)
			}
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}
