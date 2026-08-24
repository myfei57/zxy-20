package plan

import (
	"fmt"

	"bms/internal/store"
)

const (
	schedulePrefix = "plans/schedules/"
	activePrefix   = "plans/active/"
	cursorPrefix   = "plans/cursors/"
)

// activeMarker is the durable activation record of a plan version.
type activeMarker struct {
	Version int  `json:"version"`
	Active  bool `json:"active"`
}

// cursorRecord is the durable distribution cursor of a plan.
type cursorRecord struct {
	Cursor int `json:"cursor"`
}

// PlanStore persists schedules, activation markers and distribution cursors.
type PlanStore struct {
	base store.Store
}

// NewStore wraps a durable store for plan records.
func NewStore(base store.Store) *PlanStore {
	return &PlanStore{base: base}
}

// Save persists one schedule plan.
func (s *PlanStore) Save(p SchedulePlan) error {
	return store.WriteJSON(s.base, schedulePrefix+p.ID+".json", p)
}

// Get loads one schedule plan by id.
func (s *PlanStore) Get(id string) (SchedulePlan, error) {
	p, err := store.ReadJSON[SchedulePlan](s.base, schedulePrefix+id+".json")
	if err != nil {
		return SchedulePlan{}, fmt.Errorf("load plan %s: %w", id, err)
	}
	return p, nil
}

// List returns every schedule plan sorted by record name.
func (s *PlanStore) List() ([]SchedulePlan, error) {
	names, err := s.base.List(schedulePrefix)
	if err != nil {
		return nil, err
	}
	plans := make([]SchedulePlan, 0, len(names))
	for _, name := range names {
		p, err := store.ReadJSON[SchedulePlan](s.base, name)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

// MarkActive durably records that a plan version became active.
func (s *PlanStore) MarkActive(id string, version int) error {
	return store.WriteJSON(s.base, activePrefix+id+".json", activeMarker{Version: version, Active: true})
}

// Active reports the durable activation state of a plan.
func (s *PlanStore) Active(id string) (int, bool, error) {
	marker, err := store.ReadJSON[activeMarker](s.base, activePrefix+id+".json")
	if err != nil {
		if store.IsNotFound(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return marker.Version, marker.Active, nil
}

// Cursor returns the durable distribution cursor of a plan.
func (s *PlanStore) Cursor(id string) (int, error) {
	record, err := store.ReadJSON[cursorRecord](s.base, cursorPrefix+id+".json")
	if err != nil {
		if store.IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return record.Cursor, nil
}

// SetCursor durably writes the distribution cursor of a plan.
func (s *PlanStore) SetCursor(id string, value int) error {
	return store.WriteJSON(s.base, cursorPrefix+id+".json", cursorRecord{Cursor: value})
}

// AdvanceCursor increments the durable distribution cursor and returns it.
func (s *PlanStore) AdvanceCursor(id string) (int, error) {
	current, err := s.Cursor(id)
	if err != nil {
		return 0, err
	}
	next := current + 1
	if err := s.SetCursor(id, next); err != nil {
		return 0, err
	}
	return next, nil
}
