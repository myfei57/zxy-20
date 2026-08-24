package verifycase

import (
	"testing"

	"bms/internal/alert"
	"bms/internal/clock"
	"bms/internal/meter"
	"bms/internal/quota"
)

// TestMeterCursorAfterValueDurable verifies the meter cursor only moves after
// the reading itself is durably stored.
func TestMeterCursorAfterValueDurable(t *testing.T) {
	inner := newMemStore()
	failing := &failStore{inner: inner, marker: "meter/readings/"}
	clk := clock.SystemClock{}

	readings := meter.NewReadingStore(failing)
	cursor := meter.NewCursorStore(inner, readings)
	quotaSvc := quota.NewService(quota.NewStore(inner))
	alertSvc := alert.NewService(alert.NewRuleStore(inner), alert.NewEventStore(inner), clk)
	meterSvc := meter.NewService(
		meter.NewMeterStore(inner),
		readings,
		cursor,
		meter.NewWindowStore(inner, clk),
		meter.NewSummaryStore(inner),
		quotaSvc.Checker(),
		alertSvc.Checker(),
		clk,
	)
	m, err := meterSvc.CreateMeter("room-1", "电表 01")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := meterSvc.Collect(m.ID, 12.5); err == nil {
		t.Fatal("collect must fail when the reading cannot be stored")
	}
	cur, err := cursor.Current(m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if cur != 0 {
		t.Fatalf("cursor advanced to %d although the reading never landed", cur)
	}
}
