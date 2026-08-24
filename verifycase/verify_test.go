package verifycase

import (
	"testing"

	"bms/internal/alert"
	"bms/internal/clock"
	"bms/internal/meter"
	"bms/internal/quota"
)

// TestMeterWindowAfterSummaryDurable verifies a window is only marked closed
// after its summary is durably stored.
func TestMeterWindowAfterSummaryDurable(t *testing.T) {
	inner := newMemStore()
	failing := &failStore{inner: inner, marker: "meter/summaries/"}
	clk := clock.SystemClock{}

	readings := meter.NewReadingStore(inner)
	windowStore := meter.NewWindowStore(inner, clk)
	quotaSvc := quota.NewService(quota.NewStore(inner))
	alertSvc := alert.NewService(alert.NewRuleStore(inner), alert.NewEventStore(inner), clk)
	meterSvc := meter.NewService(
		meter.NewMeterStore(inner),
		readings,
		meter.NewCursorStore(inner, readings),
		windowStore,
		meter.NewSummaryStore(failing),
		quotaSvc.Checker(),
		alertSvc.Checker(),
		clk,
	)
	m, err := meterSvc.CreateMeter("room-1", "电表 01")
	if err != nil {
		t.Fatal(err)
	}
	w, err := meterSvc.OpenWindow(m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := meterSvc.CloseAfterSummary(w.ID); err == nil {
		t.Fatal("window close must fail when the summary cannot be stored")
	}
	closed, err := meterSvc.WindowClosed(w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if closed {
		t.Fatal("window marked closed although its summary never landed")
	}
}
