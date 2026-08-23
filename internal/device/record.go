package device

import (
	"fmt"
	"time"

	"bms/internal/store"
)

const commandPrefix = "device/commands/"

// CommandStore persists the durable command trace for devices.
type CommandStore struct {
	base store.Store
}

// NewCommandStore wraps a durable store for command records.
func NewCommandStore(base store.Store) *CommandStore {
	return &CommandStore{base: base}
}

// Save persists one command record.
func (s *CommandStore) Save(rec CommandRecord) error {
	return store.WriteJSON(s.base, commandPrefix+rec.ID+".json", rec)
}

// Get loads one command record by id.
func (s *CommandStore) Get(id string) (CommandRecord, error) {
	rec, err := store.ReadJSON[CommandRecord](s.base, commandPrefix+id+".json")
	if err != nil {
		return CommandRecord{}, fmt.Errorf("load command %s: %w", id, err)
	}
	return rec, nil
}

// ListByDevice returns the command trace of one device sorted by record name.
func (s *CommandStore) ListByDevice(deviceID string) ([]CommandRecord, error) {
	names, err := s.base.List(commandPrefix)
	if err != nil {
		return nil, err
	}
	records := make([]CommandRecord, 0, len(names))
	for _, name := range names {
		rec, err := store.ReadJSON[CommandRecord](s.base, name)
		if err != nil {
			return nil, err
		}
		if rec.DeviceID == deviceID {
			records = append(records, rec)
		}
	}
	return records, nil
}

// ListReplayable returns commands that were sent but may not have been
// acknowledged yet; the replay logic must still consult the ack marker.
func (s *CommandStore) ListReplayable(deviceID string) ([]CommandRecord, error) {
	records, err := s.ListByDevice(deviceID)
	if err != nil {
		return nil, err
	}
	out := make([]CommandRecord, 0, len(records))
	for _, rec := range records {
		if rec.Status == StatusSent || rec.Status == StatusAcked {
			out = append(out, rec)
		}
	}
	return out, nil
}

// Ack durably marks a command as acknowledged.
func (s *CommandStore) Ack(deviceID, commandID string, at time.Time) (CommandRecord, error) {
	rec, err := s.Get(commandID)
	if err != nil {
		return CommandRecord{}, err
	}
	if rec.DeviceID != deviceID {
		return CommandRecord{}, fmt.Errorf("command %s does not belong to device %s", commandID, deviceID)
	}
	rec.Status = StatusAcked
	rec.AckedAt = at
	if err := s.Save(rec); err != nil {
		return CommandRecord{}, err
	}
	return rec, nil
}

// IsAcked reports whether a command already has a durable acknowledgement.
func (s *CommandStore) IsAcked(deviceID, commandID string) bool {
	rec, err := s.Get(commandID)
	if err != nil {
		return false
	}
	return rec.DeviceID == deviceID && rec.Status == StatusAcked
}
