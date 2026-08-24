package audit

import "time"

// AuditEvent is one durable audit trail record.
type AuditEvent struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	EntityID   string    `json:"entity_id"`
	Result     string    `json:"result"`
	Message    string    `json:"message"`
	OccurredAt time.Time `json:"occurred_at"`
}
