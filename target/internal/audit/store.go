package audit

import "bms/internal/store"

const eventPrefix = "audit/events/"

// Store persists audit events durably.
type Store struct {
	base store.Store
}

// NewStore wraps a durable store for audit events.
func NewStore(base store.Store) *Store {
	return &Store{base: base}
}

// Append persists one audit event.
func (s *Store) Append(event AuditEvent) error {
	return store.WriteJSON(s.base, eventPrefix+event.ID+".json", event)
}

// List returns every audit event sorted by record name.
func (s *Store) List() ([]AuditEvent, error) {
	names, err := s.base.List(eventPrefix)
	if err != nil {
		return nil, err
	}
	events := make([]AuditEvent, 0, len(names))
	for _, name := range names {
		event, err := store.ReadJSON[AuditEvent](s.base, name)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// ListByKind returns audit events of one kind.
func (s *Store) ListByKind(kind string) ([]AuditEvent, error) {
	events, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]AuditEvent, 0, len(events))
	for _, event := range events {
		if event.Kind == kind {
			out = append(out, event)
		}
	}
	return out, nil
}

// Recent returns the newest events, newest first.
func (s *Store) Recent(limit int) ([]AuditEvent, error) {
	events, err := s.List()
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}
	out := make([]AuditEvent, 0, limit)
	for i := len(events) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, events[i])
	}
	return out, nil
}
