package meter

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/clock"
	"bms/internal/store"
)

const windowPrefix = "meter/windows/"

// WindowStore persists aggregation windows durably.
type WindowStore struct {
	base  store.Store
	clock clock.Clock
}

// NewWindowStore wraps a durable store for meter windows.
func NewWindowStore(base store.Store, clock clock.Clock) *WindowStore {
	return &WindowStore{base: base, clock: clock}
}

// Save persists one window record.
func (s *WindowStore) Save(w Window) error {
	return store.WriteJSON(s.base, windowPrefix+w.ID+".json", w)
}

// Get loads one window by id.
func (s *WindowStore) Get(id string) (Window, error) {
	w, err := store.ReadJSON[Window](s.base, windowPrefix+id+".json")
	if err != nil {
		return Window{}, fmt.Errorf("load window %s: %w", id, err)
	}
	return w, nil
}

// List returns the windows of a meter sorted by record name.
func (s *WindowStore) List(meterID string) ([]Window, error) {
	names, err := s.base.List(windowPrefix)
	if err != nil {
		return nil, err
	}
	windows := make([]Window, 0, len(names))
	for _, name := range names {
		w, err := store.ReadJSON[Window](s.base, name)
		if err != nil {
			return nil, err
		}
		if w.MeterID == meterID {
			windows = append(windows, w)
		}
	}
	return windows, nil
}

// MarkClosed durably marks a window as closed.
func (s *WindowStore) MarkClosed(id string) error {
	w, err := s.Get(id)
	if err != nil {
		return err
	}
	w.Closed = true
	return s.Save(w)
}

// IsClosed reports the durable close state of a window.
func (s *WindowStore) IsClosed(id string) (bool, error) {
	w, err := s.Get(id)
	if err != nil {
		return false, err
	}
	return w.Closed, nil
}

// OpenWindow creates a new open aggregation window for a meter.
func (svc *Service) OpenWindow(meterID string) (Window, error) {
	w := Window{
		ID:      uuid.NewString(),
		MeterID: meterID,
		Start:   svc.clock.Now(),
		Closed:  false,
	}
	if err := svc.windows.Save(w); err != nil {
		return Window{}, err
	}
	return w, nil
}

// CloseAfterSummary aggregates the window readings and closes the window.
func (svc *Service) CloseAfterSummary(windowID string) error {
	w, err := svc.windows.Get(windowID)
	if err != nil {
		return err
	}
	readings, err := svc.readings.List(w.MeterID)
	if err != nil {
		return err
	}
	summary := buildSummary(w, readings)
	if err := svc.windows.MarkClosed(w.ID); err != nil {
		return err
	}
	return svc.summaries.Save(summary)
}

// Windows returns the windows of a meter.
func (svc *Service) Windows(meterID string) ([]Window, error) {
	return svc.windows.List(meterID)
}
