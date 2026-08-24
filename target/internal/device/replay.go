package device

// Replayer rebuilds the command set for a reconnecting device.
type Replayer struct {
	records *CommandStore
}

// NewReplayer wires command replay over the command store.
func NewReplayer(records *CommandStore) *Replayer {
	return &Replayer{records: records}
}

// Replay returns the commands that still need to be re-sent after a device
// reconnect. Commands with a durable acknowledgement are never replayed.
func (r *Replayer) Replay(deviceID string) ([]CommandRecord, error) {
	candidates, err := r.records.ListReplayable(deviceID)
	if err != nil {
		return nil, err
	}
	out := make([]CommandRecord, 0, len(candidates))
	for _, rec := range candidates {
		if r.records.IsAcked(rec.DeviceID, rec.ID) {
			continue
		}
		out = append(out, rec)
	}
	return out, nil
}
