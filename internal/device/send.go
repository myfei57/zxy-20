package device

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/clock"
)

// Sender issues commands to devices and tracks the device state.
type Sender struct {
	records *CommandStore
	state   *StateStore
	clock   clock.Clock
}

// NewSender wires command sending over the command and state stores.
func NewSender(records *CommandStore, state *StateStore, clock clock.Clock) *Sender {
	return &Sender{records: records, state: state, clock: clock}
}

// Send issues one command to a device.
func (s *Sender) Send(deviceID, command string) (*CommandRecord, error) {
	rec := &CommandRecord{
		ID:       uuid.NewString(),
		DeviceID: deviceID,
		Command:  command,
		Status:   StatusPending,
	}
	rec.Status = StatusSent
	rec.SentAt = s.clock.Now()
	if err := s.state.MarkSent(deviceID, rec.ID, rec.SentAt); err != nil {
		return nil, err
	}
	if err := s.records.Save(*rec); err != nil {
		return nil, fmt.Errorf("record command for %s: %w", deviceID, err)
	}
	return rec, nil
}
