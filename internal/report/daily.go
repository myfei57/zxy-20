package report

import (
	"bms/internal/clock"
	"bms/internal/meter"
)

// Builder assembles daily energy reports from durable meter readings.
type Builder struct {
	meters *meter.Service
}

// NewBuilder wires daily reporting over the meter component.
func NewBuilder(meters *meter.Service) *Builder {
	return &Builder{meters: meters}
}

// Build returns one daily report per room for the given day (YYYY-MM-DD).
func (b *Builder) Build(day string) ([]DailyReport, error) {
	meters, err := b.meters.Meters()
	if err != nil {
		return nil, err
	}
	reports := make([]DailyReport, 0, len(meters))
	for _, m := range meters {
		readings, err := b.meters.Readings(m.ID)
		if err != nil {
			return nil, err
		}
		report := DailyReport{RoomID: m.RoomID, Day: day}
		for _, r := range readings {
			if clock.DayKey(r.TakenAt) != day {
				continue
			}
			report.Total += r.Value
			report.Readings++
			if r.Value > report.Peak {
				report.Peak = r.Value
			}
		}
		reports = append(reports, report)
	}
	return reports, nil
}
