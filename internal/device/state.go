package device

import (
	"fmt"
	"sync"
	"time"

	"bms/internal/store"
)

const statePrefix = "device/states/"

// StateStore keeps the durable device state plus an in-memory snapshot used by
// the room pages.
type StateStore struct {
	base     store.Store
	mu       sync.RWMutex
	snapshot map[string]string
}

// NewStateStore returns a state store with an empty snapshot.
func NewStateStore(base store.Store) *StateStore {
	return &StateStore{base: base, snapshot: make(map[string]string)}
}

// CurrentState returns the live state of a device, falling back to the durable
// record when the snapshot has not been warmed.
func (s *StateStore) CurrentState(deviceID string) string {
	s.mu.RLock()
	state, ok := s.snapshot[deviceID]
	s.mu.RUnlock()
	if ok {
		return state
	}
	dev, err := store.ReadJSON[Device](s.base, statePrefix+deviceID+".json")
	if err != nil {
		return StateOff
	}
	return dev.State
}

// SetState persists the new state and warms the in-memory snapshot.
func (s *StateStore) SetState(deviceID, state string, at time.Time) error {
	dev, err := store.ReadJSON[Device](s.base, statePrefix+deviceID+".json")
	if err != nil {
		return fmt.Errorf("state for device %s: %w", deviceID, err)
	}
	dev.State = state
	dev.UpdatedAt = at
	if err := store.WriteJSON(s.base, statePrefix+deviceID+".json", dev); err != nil {
		return err
	}
	s.mu.Lock()
	s.snapshot[deviceID] = state
	s.mu.Unlock()
	return nil
}

// Init seeds the durable state record and warms the snapshot for a new device.
func (s *StateStore) Init(deviceID, state string, at time.Time) error {
	dev := Device{ID: deviceID, State: state, UpdatedAt: at}
	if err := store.WriteJSON(s.base, statePrefix+deviceID+".json", dev); err != nil {
		return err
	}
	s.mu.Lock()
	s.snapshot[deviceID] = state
	s.mu.Unlock()
	return nil
}

// MarkSent records the latest command on a device, durably and in memory.
func (s *StateStore) MarkSent(deviceID, commandID string, at time.Time) error {
	dev, err := store.ReadJSON[Device](s.base, statePrefix+deviceID+".json")
	if err != nil {
		return fmt.Errorf("state for device %s: %w", deviceID, err)
	}
	dev.LastCommandID = commandID
	dev.LastSentAt = at
	dev.UpdatedAt = at
	if err := store.WriteJSON(s.base, statePrefix+deviceID+".json", dev); err != nil {
		return err
	}
	s.mu.Lock()
	s.snapshot[deviceID] = dev.State
	s.mu.Unlock()
	return nil
}
