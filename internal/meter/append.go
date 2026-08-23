package meter

import (
	"fmt"

	"bms/internal/store"
)

const readingPrefix = "meter/readings/"

// ReadingStore persists meter readings durably.
type ReadingStore struct {
	base store.Store
}

// NewReadingStore wraps a durable store for meter readings.
func NewReadingStore(base store.Store) *ReadingStore {
	return &ReadingStore{base: base}
}

// Append persists one reading.
func (s *ReadingStore) Append(r Reading) error {
	return store.WriteJSON(s.base, readingName(r.MeterID, r.Sequence), r)
}

// List returns the readings of a meter sorted by sequence.
func (s *ReadingStore) List(meterID string) ([]Reading, error) {
	names, err := s.base.List(readingPrefix + meterID + "-")
	if err != nil {
		return nil, err
	}
	readings := make([]Reading, 0, len(names))
	for _, name := range names {
		r, err := store.ReadJSON[Reading](s.base, name)
		if err != nil {
			return nil, err
		}
		readings = append(readings, r)
	}
	return readings, nil
}

// Ingest accepts a meter reading into the room's energy record.
func (svc *Service) Ingest(meterID string, value float64) error {
	m, err := svc.meters.Get(meterID)
	if err != nil {
		return err
	}
	sequence, err := svc.cursor.Current(meterID)
	if err != nil {
		return err
	}
	reading := Reading{MeterID: meterID, Sequence: sequence + 1, Value: value, TakenAt: svc.clock.Now()}
	if err := svc.quota.Check(m.RoomID, value); err != nil {
		return err
	}
	if err := svc.readings.Append(reading); err != nil {
		return err
	}
	if err := svc.cursor.Advance(meterID, reading.Sequence); err != nil {
		return err
	}
	if err := svc.quota.Consume(m.RoomID, value); err != nil {
		return err
	}
	if err := svc.alerts.Evaluate(m.RoomID, value); err != nil {
		return err
	}
	m.LastValue = value
	m.UpdatedAt = svc.clock.Now()
	return svc.meters.Save(m)
}

func readingName(meterID string, sequence int64) string {
	return fmt.Sprintf("%s%s-%08d.json", readingPrefix, meterID, sequence)
}
