package room

import "bms/internal/device"

// View is the console representation of one room plus its live device state.
type View struct {
	Room        Room   `json:"room"`
	DeviceState string `json:"device_state"`
}

// Query builds room views for the console pages. It reads device state live
// from the device service on every request so the page always reflects the
// latest acknowledged state rather than a stale in-process snapshot.
type Query struct {
	store   *Store
	cache   *Cache
	devices *device.Service
}

// NewQuery wires room view building over the room and device services.
func NewQuery(store *Store, cache *Cache, devices *device.Service) *Query {
	return &Query{store: store, cache: cache, devices: devices}
}

// deviceState returns the live device state bound to a room, or "unbound" when
// the room has no device attached.
func (q *Query) deviceState(r Room) string {
	if r.DeviceID == "" {
		return "unbound"
	}
	return q.devices.CurrentState(r.DeviceID)
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
	return View{Room: r, DeviceState: q.deviceState(r)}, nil
}

// All returns live views for every room.
func (q *Query) All() ([]View, error) {
	rooms, err := q.store.ListRooms()
	if err != nil {
		return nil, err
	}
	views := make([]View, 0, len(rooms))
	for _, r := range rooms {
		views = append(views, View{Room: r, DeviceState: q.deviceState(r)})
	}
	return views, nil
}
