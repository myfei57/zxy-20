package plan

import (
	"bms/internal/meter"
)

// Distributor delivers schedules to rooms.
type Distributor struct {
	plans *PlanStore
	acks  *meter.AckStore
}

// NewDistributor wires plan distribution over the plan store and the durable
// acknowledgement store.
func NewDistributor(plans *PlanStore, acks *meter.AckStore) *Distributor {
	return &Distributor{plans: plans, acks: acks}
}

// Distribute delivers the plan to the given rooms.
func (d *Distributor) Distribute(planID string, roomIDs []string) error {
	for _, roomID := range roomIDs {
		if _, err := d.plans.AdvanceCursor(planID); err != nil {
			return err
		}
		if err := d.acks.Record(roomID, planID); err != nil {
			return err
		}
	}
	return nil
}

// Acked reports whether a room already acknowledged a plan.
func (d *Distributor) Acked(roomID, planID string) bool {
	return d.acks.Acked(roomID, planID)
}
