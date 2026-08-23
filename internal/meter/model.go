package meter

import "time"

// Meter is an energy meter installed in a room.
type Meter struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	Name      string    `json:"name"`
	Cursor    int64     `json:"cursor"`
	LastValue float64   `json:"last_value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Reading is one durable meter sample.
type Reading struct {
	MeterID  string    `json:"meter_id"`
	Sequence int64     `json:"sequence"`
	Value    float64   `json:"value"`
	TakenAt  time.Time `json:"taken_at"`
}

// Window is an hourly aggregation bucket for a meter.
type Window struct {
	ID      string    `json:"id"`
	MeterID string    `json:"meter_id"`
	Start   time.Time `json:"start"`
	Closed  bool      `json:"closed"`
	Peak    float64   `json:"peak"`
	Count   int       `json:"count"`
}

// WindowSummary is the durable aggregate produced when a window closes.
type WindowSummary struct {
	WindowID string    `json:"window_id"`
	MeterID  string    `json:"meter_id"`
	Total    float64   `json:"total"`
	Peak     float64   `json:"peak"`
	Count    int       `json:"count"`
	BuiltAt  time.Time `json:"built_at"`
}

// RoomAck is the durable acknowledgement of a schedule delivered to a room.
type RoomAck struct {
	RoomID  string    `json:"room_id"`
	PlanID  string    `json:"plan_id"`
	AckedAt time.Time `json:"acked_at"`
}
