package meter

import (
	"bms/internal/store"
)

const summaryPrefix = "meter/summaries/"

// SummaryStore persists window summaries durably.
type SummaryStore struct {
	base store.Store
}

// NewSummaryStore wraps a durable store for window summaries.
func NewSummaryStore(base store.Store) *SummaryStore {
	return &SummaryStore{base: base}
}

// Save persists one window summary.
func (s *SummaryStore) Save(summary WindowSummary) error {
	return store.WriteJSON(s.base, summaryPrefix+summary.WindowID+".json", summary)
}

// buildSummary aggregates the readings of a window into a durable summary.
func buildSummary(w Window, readings []Reading) WindowSummary {
	summary := WindowSummary{WindowID: w.ID, MeterID: w.MeterID, BuiltAt: w.Start}
	for _, r := range readings {
		summary.Total += r.Value
		summary.Count++
		if r.Value > summary.Peak {
			summary.Peak = r.Value
		}
	}
	return summary
}
