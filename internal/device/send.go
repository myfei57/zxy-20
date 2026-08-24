package device

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/clock"
)

// Sender issues commands to devices, durably recording each command before the
// device is marked as having received it.
type Sender struct {
	records *CommandStore
	state   *StateStore
	clock   clock.Clock
}

// NewSender wires command sending over the command and state stores.
func NewSender(records *CommandStore, state *StateStore, clock clock.Clock) *Sender {
	return &Sender{records: records, state: state, clock: clock}
}

// Send persists a command record, marks it sent, and only then updates the
// device state so a failed record never looks like a delivered command.
func (s *Sender) Send(deviceID, command string) (*CommandRecord, error) {
	rec := &CommandRecord{
		ID:       uuid.NewString(),
		DeviceID: deviceID,
		Command:  command,
		Status:   StatusPending,
	}
	if err := s.records.Save(*rec); err != nil {
		return nil, fmt.Errorf("record command for %s: %w", deviceID, err)
	}
	rec.Status = StatusSent
	rec.SentAt = s.clock.Now()
	if err := s.records.Save(*rec); err != nil {
		return nil, fmt.Errorf("mark command sent for %s: %w", deviceID, err)
	}
	if err := s.state.MarkSent(deviceID, rec.ID, rec.SentAt); err != nil {
		return nil, err
	}
	return rec, nil
}
