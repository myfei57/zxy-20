package verifycase

import (
	"testing"

	"bms/internal/audit"
	"bms/internal/clock"
	"bms/internal/device"
	"bms/internal/plan"
	"bms/internal/room"
)

// TestAuditAfterPlanDurable verifies the audit trail only reports a successful
// plan switch after the plan itself is durably stored.
func TestAuditAfterPlanDurable(t *testing.T) {
	inner := newMemStore()
	failing := &failStore{inner: inner, marker: "plans/schedules/"}
	clk := clock.SystemClock{}

	roomStore := room.NewStore(inner)
	roomCache := room.NewCache()
	deviceSvc := device.NewService(device.NewDeviceStore(inner), device.NewCommandStore(inner), device.NewStateStore(inner), clk)
	roomSvc := room.NewService(roomStore, roomCache, deviceSvc, clk)
	auditSvc := audit.NewService(audit.NewStore(inner), clk)
	switcher := plan.NewSwitcher(plan.NewStore(failing), roomSvc.Binder(), auditSvc.Recorder(), clk)

	p := plan.SchedulePlan{ID: "plan-1", Name: "夏季计划", BuildingID: "b1", Version: 3}
	if err := switcher.Switch(p); err == nil {
		t.Fatal("switch must fail when the plan write fails")
	}
	events, err := auditSvc.Events()
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		if event.Kind == "plan.switch" && event.Result == "success" {
			t.Fatal("audit reported a successful switch although the plan never landed")
		}
	}
}
