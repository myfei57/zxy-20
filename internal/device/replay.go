package device

// Replayer rebuilds the command set for a reconnecting device.
type Replayer struct {
	records *CommandStore
}

// NewReplayer wires command replay over the command store.
func NewReplayer(records *CommandStore) *Replayer {
	return &Replayer{records: records}
}

// Replay returns the commands to re-send after a device reconnect.
func (r *Replayer) Replay(deviceID string) ([]CommandRecord, error) {
	candidates, err := r.records.ListReplayable(deviceID)
	if err != nil {
		return nil, err
	}
	return candidates, nil
}
