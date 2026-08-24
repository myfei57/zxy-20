package plan

import (
	"bms/internal/meter"
)

// Distributor pushes a schedule to rooms and advances the plan cursor only
// after every room acknowledgement is durable.
type Distributor struct {
	plans *PlanStore
	acks  *meter.AckStore
}

// NewDistributor wires plan distribution over the plan store and the durable
// acknowledgement store.
func NewDistributor(plans *PlanStore, acks *meter.AckStore) *Distributor {
	return &Distributor{plans: plans, acks: acks}
}

// Distribute records one durable acknowledgement per room and then advances
// the plan cursor, so a failed ack never leaves a skipped room behind.
func (d *Distributor) Distribute(planID string, roomIDs []string) error {
	for _, roomID := range roomIDs {
		if err := d.acks.Record(roomID, planID); err != nil {
			return err
		}
		if _, err := d.plans.AdvanceCursor(planID); err != nil {
			return err
		}
	}
	return nil
}

// Acked reports whether a room already acknowledged a plan.
func (d *Distributor) Acked(roomID, planID string) bool {
	return d.acks.Acked(roomID, planID)
}
