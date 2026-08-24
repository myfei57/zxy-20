package seed

import "time"

// today returns the local day bucket used by the demo report.
func today() string {
	return time.Now().Format("2006-01-02")
}
