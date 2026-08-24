package clock

import "time"

// Clock abstracts time so services stay deterministic in tests.
type Clock interface {
	Now() time.Time
}

// SystemClock returns wall-clock time.
type SystemClock struct{}

// NewSystemClock returns a ready-to-use system clock.
func NewSystemClock() SystemClock {
	return SystemClock{}
}

// Now implements Clock.
func (SystemClock) Now() time.Time {
	return time.Now()
}

// DayKey formats a timestamp into the daily bucket used by reports.
func DayKey(t time.Time) string {
	return t.Format("2006-01-02")
}
