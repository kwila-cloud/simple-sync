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
func NewAclService(storage storage.Storage) (*AclService, error) {
	service := &AclService{
		storage: storage,
	}
	err := service.loadRules()
	if err != nil {
		return nil, err
	}
	return service, nil
}

// loadRules loads ACL rules from storage
func (s *AclService) loadRules() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	rules, err := s.storage.GetAclRules()
	if err != nil {
		log.Printf("Failed to load ACL rules: %v", err)
		return err
	}

	s.rules = rules
	return nil
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

	// Sort by specificity (hierarchical: item > user > action > existing order)
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
		return i > j
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

// AddRule adds a new ACL rule
func (s *AclService) AddRule(rule models.AclRule) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := s.storage.CreateAclRule(&rule)
	if err != nil {
		return err
	}

	s.rules = append(s.rules, rule)
	return nil
}

// Calculates the specificity score for a pattern (wildcards worth 0.5)
func calculateSpecificity(pattern string) float64 {
	score := float64(len(pattern))
	if strings.HasSuffix(pattern, "*") {
		score -= 0.5
	}
	return score
}
