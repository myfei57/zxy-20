package alert

import "time"

// ThresholdRule is a durable threshold that is only applied once activated.
type ThresholdRule struct {
	ID      string  `json:"id"`
	RoomID  string  `json:"room_id"`
	Kind    string  `json:"kind"`
	Value   float64 `json:"value"`
	Version int     `json:"version"`
	Active  bool    `json:"active"`
}

// AlertEvent is one triggered alert.
type AlertEvent struct {
	ID          string    `json:"id"`
	RuleID      string    `json:"rule_id"`
	RoomID      string    `json:"room_id"`
	Value       float64   `json:"value"`
	TriggeredAt time.Time `json:"triggered_at"`
}
