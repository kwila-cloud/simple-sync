package performance

import (
	"strconv"
	"time"

	"simple-sync/src/models"
)

// GenerateEvents creates n events for performance testing
func GenerateEvents(n int) []models.Event {
	events := make([]models.Event, 0, n)
	for i := 0; i < n; i++ {
		uuid := "event-" + strconv.Itoa(i)
		timestamp := uint64(time.Now().Unix())
		e := models.Event{
			UUID:      uuid,
			Timestamp: timestamp,
			User:      "test-user",
			Item:      "item-" + strconv.Itoa(i%100),
			Action:    "test.action",
			Payload:   "{}",
		}
		events = append(events, e)
	}
	return events
}
