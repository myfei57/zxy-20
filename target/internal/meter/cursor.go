package meter

import (
	"bms/internal/store"
)

const cursorPrefix = "meter/cursors/"

type cursorRecord struct {
	Cursor int64 `json:"cursor"`
}

// CursorStore tracks the last durably appended reading sequence per meter.
type CursorStore struct {
	base     store.Store
	readings *ReadingStore
}

// NewCursorStore wires cursor tracking over the durable store and the reading
// store the cursor is ordered against.
func NewCursorStore(base store.Store, readings *ReadingStore) *CursorStore {
	return &CursorStore{base: base, readings: readings}
}

// Current returns the durable cursor of a meter.
func (s *CursorStore) Current(meterID string) (int64, error) {
	record, err := store.ReadJSON[cursorRecord](s.base, cursorPrefix+meterID+".json")
	if err != nil {
		if store.IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return record.Cursor, nil
}

// Advance durably writes a new cursor value.
func (s *CursorStore) Advance(meterID string, sequence int64) error {
	return store.WriteJSON(s.base, cursorPrefix+meterID+".json", cursorRecord{Cursor: sequence})
}

// AdvanceAfterValue persists a reading and advances the meter cursor.
func (s *CursorStore) AdvanceAfterValue(r Reading) error {
	if err := s.Advance(r.MeterID, r.Sequence); err != nil {
		return err
	}
	return s.readings.Append(r)
}

// Collect reads one meter sample and advances the meter cursor.
func (svc *Service) Collect(meterID string, value float64) (Reading, error) {
	m, err := svc.meters.Get(meterID)
	if err != nil {
		return Reading{}, err
	}
	sequence, err := svc.cursor.Current(meterID)
	if err != nil {
		return Reading{}, err
	}
	reading := Reading{MeterID: meterID, Sequence: sequence + 1, Value: value, TakenAt: svc.clock.Now()}
	if err := svc.cursor.AdvanceAfterValue(reading); err != nil {
		return Reading{}, err
	}
	if err := svc.alerts.Evaluate(m.RoomID, value); err != nil {
		return Reading{}, err
	}
	m.LastValue = value
	m.UpdatedAt = svc.clock.Now()
	if err := svc.meters.Save(m); err != nil {
		return Reading{}, err
	}
	return reading, nil
}
