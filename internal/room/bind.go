package room

import (
	"fmt"

	"bms/internal/clock"
)

// Binder binds rooms to a running schedule plan.
type Binder struct {
	store *Store
	cache *Cache
	clock clock.Clock
}

// NewBinder wires room binding over a room store and cache.
func NewBinder(store *Store, cache *Cache, clock clock.Clock) *Binder {
	return &Binder{store: store, cache: cache, clock: clock}
}

// BindToPlan durably records that a room now follows the given plan.
func (b *Binder) BindToPlan(roomID, planID string) error {
	r, ok := b.cache.Get(roomID)
	if !ok {
		var err error
		r, err = b.store.GetRoom(roomID)
		if err != nil {
			return fmt.Errorf("bind room %s: %w", roomID, err)
		}
	}
	updated := r
	updated.BoundPlanID = planID
	updated.UpdatedAt = b.clock.Now()
	if err := b.store.SaveRoom(updated); err != nil {
		return fmt.Errorf("persist room binding %s: %w", roomID, err)
	}
	b.cache.Set(updated)
	return nil
}

// BindRoomsToPlan binds every room in a building to the plan.
func (b *Binder) BindRoomsToPlan(buildingID, planID string) error {
	rooms, err := b.store.ListByBuilding(buildingID)
	if err != nil {
		return err
	}
	for _, r := range rooms {
		if err := b.BindToPlan(r.ID, planID); err != nil {
			return err
		}
	}
	return nil
}
