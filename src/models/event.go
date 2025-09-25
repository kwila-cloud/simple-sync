package models

// Event represents a timestamped event in the system
type Event struct {
	UUID      string `json:"uuid"`
	Timestamp uint64 `json:"timestamp"`
	User      string `json:"user"`
	Item      string `json:"item"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
}
