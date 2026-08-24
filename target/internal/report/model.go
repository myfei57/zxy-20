package report

// DailyReport is the per-room energy aggregate for one day.
type DailyReport struct {
	RoomID   string  `json:"room_id"`
	Day      string  `json:"day"`
	Total    float64 `json:"total"`
	Peak     float64 `json:"peak"`
	Readings int     `json:"readings"`
}
