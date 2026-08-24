package room

import "bms/internal/device"

// View is the console representation of one room plus its live device state.
type View struct {
	Room        Room   `json:"room"`
	DeviceState string `json:"device_state"`
}

// Query builds room views for the console pages.
type Query struct {
	store     *Store
	cache     *Cache
	devices   *device.Service
	snapshots map[string]string
}

// NewQuery wires room view building over the room and device services.
func NewQuery(store *Store, cache *Cache, devices *device.Service) *Query {
	return &Query{store: store, cache: cache, devices: devices, snapshots: make(map[string]string)}
}

// ByID returns the console view of one room.
func (q *Query) ByID(roomID string) (View, error) {
	r, ok := q.cache.Get(roomID)
	if !ok {
		var err error
		r, err = q.store.GetRoom(roomID)
		if err != nil {
			return View{}, err
		}
	}
	state := "unbound"
	if r.DeviceID != "" {
		state, ok = q.snapshots[r.DeviceID]
		if !ok {
			state = q.devices.CurrentState(r.DeviceID)
			q.snapshots[r.DeviceID] = state
		}
	}
	return View{Room: r, DeviceState: state}, nil
}

// All returns live views for every room.
func (q *Query) All() ([]View, error) {
	rooms, err := q.store.ListRooms()
	if err != nil {
		return nil, err
	}
	views := make([]View, 0, len(rooms))
	for _, r := range rooms {
		state := "unbound"
		if r.DeviceID != "" {
			state = q.devices.CurrentState(r.DeviceID)
		}
		views = append(views, View{Room: r, DeviceState: state})
	}
	return views, nil
}
