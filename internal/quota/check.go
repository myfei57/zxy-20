package quota

import (
	"errors"
	"fmt"

	"bms/internal/store"
)

// ErrQuotaExceeded is returned when a reading would push a room past its
// energy budget.
var ErrQuotaExceeded = errors.New("energy quota exceeded")

// Checker gates meter ingestion against the durable energy quota.
type Checker struct {
	store *Store
}

// NewChecker wires quota gating over the quota store.
func NewChecker(store *Store) *Checker {
	return &Checker{store: store}
}

// Check returns an error when ingesting value would exceed the room quota.
func (c *Checker) Check(roomID string, value float64) error {
	q, err := c.store.GetByRoom(roomID)
	if err != nil {
		if store.IsNotFound(err) {
			return nil
		}
		return err
	}
	if q.Used+value > q.Limit {
		return fmt.Errorf("%w: used %.2f plus %.2f exceeds limit %.2f", ErrQuotaExceeded, q.Used, value, q.Limit)
	}
	return nil
}

// Consume durably adds a value to the used energy of a room.
func (c *Checker) Consume(roomID string, value float64) error {
	q, err := c.store.GetByRoom(roomID)
	if err != nil {
		if store.IsNotFound(err) {
			return nil
		}
		return err
	}
	q.Used += value
	q.Version++
	return c.store.Save(q)
}
