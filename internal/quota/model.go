package quota

// Quota is the durable energy budget of a room.
type Quota struct {
	ID      string  `json:"id"`
	Scope   string  `json:"scope"`
	RoomID  string  `json:"room_id"`
	Limit   float64 `json:"limit"`
	Used    float64 `json:"used"`
	Version int     `json:"version"`
}
