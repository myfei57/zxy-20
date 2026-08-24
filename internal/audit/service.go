package audit

import (
	"bms/internal/clock"
)

// Service is the audit component entry point.
type Service struct {
	store    *Store
	recorder *Recorder
}

// NewService wires the audit component over its store.
func NewService(store *Store, clock clock.Clock) *Service {
	return &Service{store: store, recorder: NewRecorder(store, clock)}
}

// Recorder exposes the durable recorder to other components.
func (svc *Service) Recorder() *Recorder {
	return svc.recorder
}

// Events returns every audit event.
func (svc *Service) Events() ([]AuditEvent, error) {
	return svc.store.List()
}

// EventsByKind returns audit events of one kind.
func (svc *Service) EventsByKind(kind string) ([]AuditEvent, error) {
	return svc.store.ListByKind(kind)
}

// Recent returns the newest audit events.
func (svc *Service) Recent(limit int) ([]AuditEvent, error) {
	return svc.store.Recent(limit)
}
