package ns

// Building is a top-level namespace owned by the platform.
type Building struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Zone groups rooms inside a building.
type Zone struct {
	ID         string `json:"id"`
	BuildingID string `json:"building_id"`
	Name       string `json:"name"`
	Floor      string `json:"floor"`
}
