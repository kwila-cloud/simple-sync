package services

import (
	"log"
	"sort"
	"strings"
	"sync"

	"simple-sync/src/models"
	"simple-sync/src/storage"
)

// AclService handles access control logic
type AclService struct {
	storage storage.Storage
	rules   []models.AclRule
	mutex   sync.RWMutex
}

// NewAclService creates a new ACL service
func NewAclService(storage storage.Storage) *AclService {
	service := &AclService{
		storage: storage,
	}
	service.loadRules()
	return service
}

// loadRules loads ACL rules from storage
func (s *AclService) loadRules() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	events, err := s.storage.LoadEvents()
	if err != nil {
		log.Printf("Failed to load events for ACL: %v", err)
		return
	}

	var rules []models.AclRule
	for _, event := range events {
		if event.IsAclEvent() {
			rule, err := event.ToAclRule()
			if err != nil {
				log.Printf("Failed to parse ACL rule: %v", err)
				continue
			}
			rules = append(rules, *rule)
		}
	}

	s.rules = rules
}

// calculateSpecificity calculates the specificity score for a pattern (wildcards worth 0.5)
func calculateSpecificity(pattern string) float64 {
	score := float64(len(pattern))
	if strings.HasSuffix(pattern, "*") {
		score -= 0.5
	}
	return score
}

// isValidPattern checks if a pattern has valid wildcard usage (at most one at the end)
func isValidPattern(pattern string) bool {
	if pattern == "" {
		return false
	}
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return !strings.Contains(prefix, "*")
	}
	return !strings.Contains(pattern, "*")
}

// CheckPermission checks if a user has permission for an action on an item
func (s *AclService) CheckPermission(user, item, action string) bool {
	// Root user bypass
	if user == ".root" {
		log.Printf("ACL: Root user %s bypass for %s on %s", user, action, item)
		return true
	}

	s.mutex.RLock()
	rules := s.rules
	s.mutex.RUnlock()

	// Find applicable rules
	var applicableRules []models.AclRule
	for _, rule := range rules {
		if s.matches(rule.User, user) && s.matches(rule.Item, item) && s.matches(rule.Action, action) {
			applicableRules = append(applicableRules, rule)
		}
	}

	if len(applicableRules) == 0 {
		log.Printf("ACL: Deny by default for user=%s, item=%s, action=%s", user, item, action)
		return false // Deny by default
	}

	// Sort by specificity and timestamp (hierarchical: item > user > action > timestamp)
	sort.Slice(applicableRules, func(i, j int) bool {
		itemI := calculateSpecificity(applicableRules[i].Item)
		itemJ := calculateSpecificity(applicableRules[j].Item)
		if itemI != itemJ {
			return itemI > itemJ
		}
		userI := calculateSpecificity(applicableRules[i].User)
		userJ := calculateSpecificity(applicableRules[j].User)
		if userI != userJ {
			return userI > userJ
		}
		actionI := calculateSpecificity(applicableRules[i].Action)
		actionJ := calculateSpecificity(applicableRules[j].Action)
		if actionI != actionJ {
			return actionI > actionJ
		}
		// Fallback to using the most recent rule
		return applicableRules[i].Timestamp > applicableRules[j].Timestamp
	})

	// The first (highest specificity/latest) determines
	allowed := applicableRules[0].Type == "allow"
	log.Printf("ACL: Decision for user=%s, item=%s, action=%s: %v (rule: %s)", user, item, action, allowed, applicableRules[0].Type)
	return allowed
}

// matches checks if pattern matches value (supports wildcards)
func (s *AclService) matches(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(value, prefix)
	}
	return pattern == value
}

// ValidateAclEvent validates if an ACL event can be stored
func (s *AclService) ValidateAclEvent(event *models.Event) bool {
	if !event.IsAclEvent() {
		return false // Not ACL, reject
	}

	// Validate action must be .acl.allow or .acl.deny
	if event.Action != ".acl.allow" && event.Action != ".acl.deny" {
		return false
	}

	// Validate payload can be parsed and fields are not empty
	rule, err := event.ToAclRule()
	if err != nil {
		return false
	}

	// Validate rule patterns
	if !isValidPattern(rule.User) || !isValidPattern(rule.Item) || !isValidPattern(rule.Action) {
		return false
	}

	// Check if the user can set this ACL rule
	return s.CheckPermission(event.User, ".acl", event.Action)
}

// AddRule adds a new ACL rule (called after event is stored)
func (s *AclService) AddRule(rule models.AclRule) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.rules = append(s.rules, rule)
}

// RefreshRules reloads rules from storage
func (s *AclService) RefreshRules() {
	s.loadRules()
}
