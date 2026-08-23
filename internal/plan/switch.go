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
//
// The plan record and the room bindings are persisted before the activation
// marker is written: the marker is the only thing the console reads as the
// "active" signal, so it must not land until every room is durably bound to
// the new plan. A restart at any point then either leaves the rooms on the
// previous plan with the marker still absent (the page honestly shows the
// plan as a draft) or leaves them fully on the new plan — it can never show
// the plan as active while rooms still follow the old one.
func (sw *Switcher) Switch(p SchedulePlan) error {
	p.Active = true
	if err := sw.plans.Save(p); err != nil {
		return fmt.Errorf("persist plan %s: %w", p.ID, err)
	}
	if err := sw.rooms.BindRoomsToPlan(p.BuildingID, p.ID); err != nil {
		return err
	}
	if err := sw.plans.MarkActive(p.ID, p.Version); err != nil {
		return err
	}
	return sw.audit.Record("plan.switch", p.ID, "success", fmt.Sprintf("version %d active", p.Version))
}
