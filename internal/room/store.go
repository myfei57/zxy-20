package room

import (
	"fmt"

	"bms/internal/store"
)

const roomPrefix = "rooms/"

// Store persists room records durably.
type Store struct {
	base store.Store
}

// NewStore wraps a durable store for room records.
func NewStore(base store.Store) *Store {
	return &Store{base: base}
}

// SaveRoom persists one room record.
func (s *Store) SaveRoom(r Room) error {
	return store.WriteJSON(s.base, roomPrefix+r.ID+".json", r)
}

// GetRoom loads one room by id.
func (s *Store) GetRoom(id string) (Room, error) {
	r, err := store.ReadJSON[Room](s.base, roomPrefix+id+".json")
	if err != nil {
		return Room{}, fmt.Errorf("load room %s: %w", id, err)
	}
	return r, nil
}

// ListRooms returns every room sorted by record name.
func (s *Store) ListRooms() ([]Room, error) {
	names, err := s.base.List(roomPrefix)
	if err != nil {
		return nil, err
	}
	rooms := make([]Room, 0, len(names))
	for _, name := range names {
		r, err := store.ReadJSON[Room](s.base, name)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

// ListByBuilding returns rooms inside a building.
func (s *Store) ListByBuilding(buildingID string) ([]Room, error) {
	rooms, err := s.ListRooms()
	if err != nil {
		return nil, err
	}
	out := make([]Room, 0, len(rooms))
	for _, r := range rooms {
		if r.BuildingID == buildingID {
			out = append(out, r)
		}
	}
	return out, nil
}
