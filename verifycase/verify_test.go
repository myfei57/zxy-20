package verifycase

import (
	"testing"

	"bms/internal/audit"
	"bms/internal/clock"
	"bms/internal/device"
	"bms/internal/plan"
	"bms/internal/room"
)

// TestScheduleActiveAfterPlanDurable verifies a plan only becomes active after
// the plan bytes are durably stored.
func TestScheduleActiveAfterPlanDurable(t *testing.T) {
	inner := newMemStore()
	failing := &failStore{inner: inner, marker: "plans/schedules/"}
	clk := clock.SystemClock{}

	roomStore := room.NewStore(inner)
	roomCache := room.NewCache()
	deviceSvc := device.NewService(device.NewDeviceStore(inner), device.NewCommandStore(inner), device.NewStateStore(inner), clk)
	roomSvc := room.NewService(roomStore, roomCache, deviceSvc, clk)

	r := room.Room{ID: "room-1", BuildingID: "b1", Name: "三层会客室", Setpoint: 24}
	if err := roomStore.SaveRoom(r); err != nil {
		t.Fatal(err)
	}
	roomCache.Set(r)

	auditSvc := audit.NewService(audit.NewStore(inner), clk)
	plans := plan.NewStore(failing)
	switcher := plan.NewSwitcher(plans, roomSvc.Binder(), auditSvc.Recorder(), clk)

	p := plan.SchedulePlan{ID: "plan-1", Name: "夏季计划", BuildingID: "b1", Version: 2}
	if err := switcher.Switch(p); err == nil {
		t.Fatal("switch must fail when the plan write fails")
	}
	if _, active, err := plans.Active("plan-1"); err != nil {
		t.Fatal(err)
	} else if active {
		t.Fatal("plan marked active although its plan never landed")
	}
	got, err := roomStore.GetRoom("room-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.BoundPlanID != "" {
		t.Fatalf("room bound to plan %s that never landed", got.BoundPlanID)
	}
}
