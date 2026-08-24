package meter

import (
	"fmt"

	"bms/internal/clock"
	"bms/internal/store"
)

const ackPrefix = "meter/acks/"

// AckStore persists schedule delivery acknowledgements durably.
type AckStore struct {
	base  store.Store
	clock clock.Clock
}

// NewAckStore wraps a durable store for room acknowledgements.
func NewAckStore(base store.Store, clock clock.Clock) *AckStore {
	return &AckStore{base: base, clock: clock}
}

// Record persists one room acknowledgement.
func (s *AckStore) Record(roomID, planID string) error {
	ack := RoomAck{RoomID: roomID, PlanID: planID, AckedAt: s.clock.Now()}
	return store.WriteJSON(s.base, ackName(roomID, planID), ack)
}

// Acked reports whether a room already acknowledged a plan.
func (s *AckStore) Acked(roomID, planID string) bool {
	return s.base.Exists(ackName(roomID, planID))
}

func ackName(roomID, planID string) string {
	return fmt.Sprintf("%s%s-%s.json", ackPrefix, roomID, planID)
}
