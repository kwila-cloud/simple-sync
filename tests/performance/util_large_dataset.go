package performance

import (
	"strconv"

	"simple-sync/src/models"
)

// GenerateEvents creates n events for performance testing
func GenerateEvents(n int) []models.Event {
	events := make([]models.Event, 0, n)
	for i := 0; i < n; i++ {
		item := "item-" + strconv.Itoa(i%100)
		e := models.NewEvent("test-user", item, "test.action", "{}")
		events = append(events, *e)
	}
	return events
}
