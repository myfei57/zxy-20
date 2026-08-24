package verifycase

import (
	"testing"

	"bms/internal/alert"
	"bms/internal/clock"
	"bms/internal/meter"
	"bms/internal/quota"
)

// TestEnergyQuotaRejectsBeforeMeter verifies an over-quota reading is refused
// before it occupies any durable storage.
func TestEnergyQuotaRejectsBeforeMeter(t *testing.T) {
	inner := newMemStore()
	clk := clock.SystemClock{}

	quotaSvc := quota.NewService(quota.NewStore(inner))
	if _, err := quotaSvc.Create("room-1", 10); err != nil {
		t.Fatal(err)
	}
	alertSvc := alert.NewService(alert.NewRuleStore(inner), alert.NewEventStore(inner), clk)
	readings := meter.NewReadingStore(inner)
	meterSvc := meter.NewService(
		meter.NewMeterStore(inner),
		readings,
		meter.NewCursorStore(inner, readings),
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
	if err := meterSvc.Ingest(m.ID, 12); err == nil {
		t.Fatal("over-quota ingest must be rejected")
	}
	stored, err := meterSvc.Readings(m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(stored) != 0 {
		t.Fatalf("over-quota reading consumed storage (%d stored)", len(stored))
	}
}
