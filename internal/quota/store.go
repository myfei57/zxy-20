package quota

import (
	"fmt"

	"bms/internal/store"
)

const quotaPrefix = "quota/"

// Store persists energy quotas durably.
type Store struct {
	base store.Store
}

// NewStore wraps a durable store for quota records.
func NewStore(base store.Store) *Store {
	return &Store{base: base}
}

// Save persists one quota record.
func (s *Store) Save(q Quota) error {
	return store.WriteJSON(s.base, quotaPrefix+q.ID+".json", q)
}

// GetByRoom loads the quota of a room.
func (s *Store) GetByRoom(roomID string) (Quota, error) {
	quotas, err := s.List()
	if err != nil {
		return Quota{}, err
	}
	for _, q := range quotas {
		if q.RoomID == roomID {
			return q, nil
		}
	}
	return Quota{}, fmt.Errorf("quota for room %s: %w", roomID, store.ErrNotFound)
}

// List returns every quota sorted by record name.
func (s *Store) List() ([]Quota, error) {
	names, err := s.base.List(quotaPrefix)
	if err != nil {
		return nil, err
	}
	quotas := make([]Quota, 0, len(names))
	for _, name := range names {
		q, err := store.ReadJSON[Quota](s.base, name)
		if err != nil {
			return nil, err
		}
		quotas = append(quotas, q)
	}
	return quotas, nil
}
