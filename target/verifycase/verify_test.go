package verifycase

import (
	"testing"

	"bms/internal/clock"
	"bms/internal/meter"
	"bms/internal/plan"
)

// TestPlanCursorAfterAckDurable verifies the plan distribution cursor advances
// only after the room acknowledgement is durably stored.
func TestPlanCursorAfterAckDurable(t *testing.T) {
	inner := newMemStore()
	failing := &failStore{inner: inner, marker: "meter/acks/"}
	clk := clock.SystemClock{}

	plans := plan.NewStore(inner)
	p := plan.SchedulePlan{ID: "plan-1", Name: "夏季计划", BuildingID: "b1", Version: 1}
	if err := plans.Save(p); err != nil {
		t.Fatal(err)
	}
	distributor := plan.NewDistributor(plans, meter.NewAckStore(failing, clk))
	if err := distributor.Distribute("plan-1", []string{"room-1"}); err == nil {
		t.Fatal("distribution must fail when the room ack cannot be stored")
	}
	cur, err := plans.Cursor("plan-1")
	if err != nil {
		t.Fatal(err)
	}
	if cur != 0 {
		t.Fatalf("cursor advanced to %d although the room ack never landed", cur)
	}
}
