package models

// AclRule represents an access control rule
type AclRule struct {
	User   string `json:"user"`
	Item   string `json:"item"`
	Action string `json:"action"`
	Type   string `json:"type"`
}
