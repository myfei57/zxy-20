package alert

import (
	"github.com/google/uuid"

	"bms/internal/clock"
)

// Checker evaluates active threshold versions against new meter readings.
type Checker struct {
	rules *RuleStore
	store *EventStore
	clock clock.Clock
}

// NewChecker wires alert evaluation over the rule and event stores.
func NewChecker(rules *RuleStore, store *EventStore, clock clock.Clock) *Checker {
	return &Checker{rules: rules, store: store, clock: clock}
}

// Evaluate checks a reading against the active rules for its room and persists
// every triggered alert.
func (c *Checker) Evaluate(roomID string, value float64) error {
	rules, err := c.rules.ListActive()
	if err != nil {
		return err
	}
	for _, rule := range rules {
		if rule.RoomID != roomID && rule.RoomID != "*" {
			continue
		}
		if value >= rule.Value {
			event := AlertEvent{
				ID:          uuid.NewString(),
				RuleID:      rule.ID,
				RoomID:      roomID,
				Value:       value,
				TriggeredAt: c.clock.Now(),
			}
			if err := c.store.Append(event); err != nil {
				return err
			}
		}
	}
	return nil
}
