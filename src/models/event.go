package models

// Event represents a timestamped event in the system
type Event struct {
	UUID      string `json:"uuid"`
	Timestamp uint64 `json:"timestamp"`
	UserUUID  string `json:"userUuid"`
	ItemUUID  string `json:"itemUuid"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
}