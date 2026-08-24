package plan

import (
	"fmt"

	"bms/internal/audit"
	"bms/internal/clock"
	"bms/internal/room"
)

// Switcher activates schedule plans and refreshes room bindings.
type Switcher struct {
	plans *PlanStore
	rooms *room.Binder
	audit *audit.Recorder
	clock clock.Clock
}

// NewSwitcher wires plan activation over the plan store, room binder and audit
// recorder.
func NewSwitcher(plans *PlanStore, rooms *room.Binder, audit *audit.Recorder, clock clock.Clock) *Switcher {
	return &Switcher{plans: plans, rooms: rooms, audit: audit, clock: clock}
}

// Switch activates the given plan for its building.
func (sw *Switcher) Switch(p SchedulePlan) error {
	p.Active = true
	if err := sw.plans.MarkActive(p.ID, p.Version); err != nil {
		return err
	}
	if err := sw.rooms.BindRoomsToPlan(p.BuildingID, p.ID); err != nil {
		return err
	}
	if err := sw.plans.Save(p); err != nil {
		return fmt.Errorf("persist plan %s: %w", p.ID, err)
	}
	return sw.audit.Record("plan.switch", p.ID, "success", fmt.Sprintf("version %d active", p.Version))
}
