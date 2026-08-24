package plan

import "time"

// SchedulePlan is a time schedule bound to the rooms of a building.
type SchedulePlan struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	BuildingID string    `json:"building_id"`
	Version    int       `json:"version"`
	Active     bool      `json:"active"`
	Cursor     int       `json:"cursor"`
	CreatedAt  time.Time `json:"created_at"`
}
