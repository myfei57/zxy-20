package plan

import (
	"fmt"

	"bms/internal/audit"
	"bms/internal/clock"
	"bms/internal/room"
)

// Switcher activates a schedule plan only after the plan is durably stored,
// then refreshes room bindings and records the audit success.
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

// Switch makes p the active schedule. The plan bytes are durable before the
// activation marker, the room bindings and the audit success event.
func (sw *Switcher) Switch(p SchedulePlan) error {
	p.Active = true
	if err := sw.plans.Save(p); err != nil {
		return fmt.Errorf("persist plan %s: %w", p.ID, err)
	}
	if err := sw.plans.MarkActive(p.ID, p.Version); err != nil {
		return err
	}
	if err := sw.rooms.BindRoomsToPlan(p.BuildingID, p.ID); err != nil {
		return err
	}
	return sw.audit.Record("plan.switch", p.ID, "success", fmt.Sprintf("version %d active", p.Version))
}
