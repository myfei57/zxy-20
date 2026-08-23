package room

import (
	"fmt"

	"bms/internal/clock"
)

// SetpointService applies new setpoints to rooms.
type SetpointService struct {
	store *Store
	cache *Cache
	clock clock.Clock
}

// NewSetpointService wires setpoint updates over a room store and cache.
func NewSetpointService(store *Store, cache *Cache, clock clock.Clock) *SetpointService {
	return &SetpointService{store: store, cache: cache, clock: clock}
}

// Set applies a new setpoint to the room.
func (svc *SetpointService) Set(roomID string, value float64) error {
	r, ok := svc.cache.Get(roomID)
	if !ok {
		var err error
		r, err = svc.store.GetRoom(roomID)
		if err != nil {
			return fmt.Errorf("setpoint for %s: %w", roomID, err)
		}
	}
	updated := r
	updated.Setpoint = value
	updated.UpdatedAt = svc.clock.Now()
	if err := svc.store.SaveRoom(updated); err != nil {
		return fmt.Errorf("persist setpoint for %s: %w", roomID, err)
	}
	svc.cache.Set(updated)
	return nil
}
