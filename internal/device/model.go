package device

import "time"

// Device states surfaced on the console pages.
const (
	StateOff      = "off"
	StateStarting = "starting"
	StateRunning  = "running"
	StateStopped  = "stopped"
	StateFault    = "fault"
)

// Command statuses in the device command lifecycle.
const (
	StatusPending = "pending"
	StatusSent    = "sent"
	StatusAcked   = "acked"
)

// Device is a registered HVAC unit.
type Device struct {
	ID            string    `json:"id"`
	RoomID        string    `json:"room_id"`
	Name          string    `json:"name"`
	Kind          string    `json:"kind"`
	State         string    `json:"state"`
	LastCommandID string    `json:"last_command_id"`
	LastSentAt    time.Time `json:"last_sent_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CommandRecord is the durable trace of one command sent to a device.
type CommandRecord struct {
	ID       string    `json:"id"`
	DeviceID string    `json:"device_id"`
	Command  string    `json:"command"`
	Status   string    `json:"status"`
	SentAt   time.Time `json:"sent_at"`
	AckedAt  time.Time `json:"acked_at"`
}
