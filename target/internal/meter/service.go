package meter

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/alert"
	"bms/internal/clock"
	"bms/internal/quota"
)

// Service is the meter component entry point.
type Service struct {
	meters    *MeterStore
	readings  *ReadingStore
	cursor    *CursorStore
	windows   *WindowStore
	summaries *SummaryStore
	quota     *quota.Checker
	alerts    *alert.Checker
	clock     clock.Clock
}

// NewService wires the meter component over its stores and gates.
func NewService(
	meters *MeterStore,
	readings *ReadingStore,
	cursor *CursorStore,
	windows *WindowStore,
	summaries *SummaryStore,
	quota *quota.Checker,
	alerts *alert.Checker,
	clock clock.Clock,
) *Service {
	return &Service{
		meters:    meters,
		readings:  readings,
		cursor:    cursor,
		windows:   windows,
		summaries: summaries,
		quota:     quota,
		alerts:    alerts,
		clock:     clock,
	}
}

// CreateMeter registers a new meter in a room.
func (svc *Service) CreateMeter(roomID, name string) (Meter, error) {
	if name == "" {
		return Meter{}, fmt.Errorf("meter name is required")
	}
	m := Meter{ID: uuid.NewString(), RoomID: roomID, Name: name, UpdatedAt: svc.clock.Now()}
	if err := svc.meters.Save(m); err != nil {
		return Meter{}, err
	}
	return m, nil
}

// Meters lists every registered meter.
func (svc *Service) Meters() ([]Meter, error) {
	return svc.meters.List()
}

// Readings returns the durable readings of a meter.
func (svc *Service) Readings(meterID string) ([]Reading, error) {
	return svc.readings.List(meterID)
}

// WindowClosed reports the durable close state of an aggregation window.
func (svc *Service) WindowClosed(windowID string) (bool, error) {
	return svc.windows.IsClosed(windowID)
}
