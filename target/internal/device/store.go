package device

import (
	"fmt"

	"bms/internal/store"
)

const devicePrefix = "device/registry/"

// DeviceStore persists the device registry durably.
type DeviceStore struct {
	base store.Store
}

// NewDeviceStore wraps a durable store for device registry records.
func NewDeviceStore(base store.Store) *DeviceStore {
	return &DeviceStore{base: base}
}

// Save persists one device record.
func (s *DeviceStore) Save(dev Device) error {
	return store.WriteJSON(s.base, devicePrefix+dev.ID+".json", dev)
}

// Get loads one device by id.
func (s *DeviceStore) Get(id string) (Device, error) {
	dev, err := store.ReadJSON[Device](s.base, devicePrefix+id+".json")
	if err != nil {
		return Device{}, fmt.Errorf("load device %s: %w", id, err)
	}
	return dev, nil
}

// List returns every device sorted by record name.
func (s *DeviceStore) List() ([]Device, error) {
	names, err := s.base.List(devicePrefix)
	if err != nil {
		return nil, err
	}
	devices := make([]Device, 0, len(names))
	for _, name := range names {
		dev, err := store.ReadJSON[Device](s.base, name)
		if err != nil {
			return nil, err
		}
		devices = append(devices, dev)
	}
	return devices, nil
}
