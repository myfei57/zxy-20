package room

import "time"

// Mode is the HVAC operating mode of a room.
type Mode string

const (
	ModeAuto Mode = "auto"
	ModeCool Mode = "cool"
)

// Room is a temperature-controlled space with a durable setpoint.
type Room struct {
	ID          string    `json:"id"`
	BuildingID  string    `json:"building_id"`
	ZoneID      string    `json:"zone_id"`
	Name        string    `json:"name"`
	Setpoint    float64   `json:"setpoint"`
	Mode        Mode      `json:"mode"`
	BoundPlanID string    `json:"bound_plan_id"`
	DeviceID    string    `json:"device_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}
