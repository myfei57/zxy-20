package audit

import (
	"github.com/google/uuid"

	"bms/internal/clock"
)

// Recorder writes audit trail events durably.
type Recorder struct {
	store *Store
	clock clock.Clock
}

// NewRecorder wires audit recording over the audit store.
func NewRecorder(store *Store, clock clock.Clock) *Recorder {
	return &Recorder{store: store, clock: clock}
}

// Record appends one audit event.
func (r *Recorder) Record(kind, entityID, result, message string) error {
	event := AuditEvent{
		ID:         uuid.NewString(),
		Kind:       kind,
		EntityID:   entityID,
		Result:     result,
		Message:    message,
		OccurredAt: r.clock.Now(),
	}
	return r.store.Append(event)
}
