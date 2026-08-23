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
//
// The command record is persisted before the device state is marked as sent, so
// that the console invariant holds: whenever a device shows "sent" the matching
// command record has already been durably written. Failing to persist the
// record therefore fails the whole send without touching device state.
func (s *Sender) Send(deviceID, command string) (*CommandRecord, error) {
	rec := &CommandRecord{
		ID:       uuid.NewString(),
		DeviceID: deviceID,
		Command:  command,
		Status:   StatusSent,
		SentAt:   s.clock.Now(),
	}
	if err := s.records.Save(*rec); err != nil {
		return nil, fmt.Errorf("record command for %s: %w", deviceID, err)
	}
	if err := s.state.MarkSent(deviceID, rec.ID, rec.SentAt); err != nil {
		return nil, err
	}
	return rec, nil
}
