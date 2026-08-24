package meter

import (
	"fmt"

	"bms/internal/store"
)

const meterPrefix = "meter/registry/"

// MeterStore persists the meter registry durably.
type MeterStore struct {
	base store.Store
}

// NewMeterStore wraps a durable store for meter records.
func NewMeterStore(base store.Store) *MeterStore {
	return &MeterStore{base: base}
}

// Save persists one meter record.
func (s *MeterStore) Save(m Meter) error {
	return store.WriteJSON(s.base, meterPrefix+m.ID+".json", m)
}

// Get loads one meter by id.
func (s *MeterStore) Get(id string) (Meter, error) {
	m, err := store.ReadJSON[Meter](s.base, meterPrefix+id+".json")
	if err != nil {
		return Meter{}, fmt.Errorf("load meter %s: %w", id, err)
	}
	return m, nil
}

// List returns every meter sorted by record name.
func (s *MeterStore) List() ([]Meter, error) {
	names, err := s.base.List(meterPrefix)
	if err != nil {
		return nil, err
	}
	meters := make([]Meter, 0, len(names))
	for _, name := range names {
		m, err := store.ReadJSON[Meter](s.base, name)
		if err != nil {
			return nil, err
		}
		meters = append(meters, m)
	}
	return meters, nil
}
