package ns

import (
	"fmt"

	"bms/internal/store"
)

const (
	buildingPrefix = "ns/buildings/"
	zonePrefix     = "ns/zones/"
)

// Store persists namespace records on disk.
type Store struct {
	base store.Store
}

// NewStore wraps a durable store for namespace records.
func NewStore(base store.Store) *Store {
	return &Store{base: base}
}

// SaveBuilding persists a building record.
func (s *Store) SaveBuilding(building Building) error {
	return store.WriteJSON(s.base, buildingPrefix+building.ID+".json", building)
}

// GetBuilding loads one building by id.
func (s *Store) GetBuilding(id string) (Building, error) {
	building, err := store.ReadJSON[Building](s.base, buildingPrefix+id+".json")
	if err != nil {
		return Building{}, fmt.Errorf("load building %s: %w", id, err)
	}
	return building, nil
}

// ListBuildings returns every building sorted by record name.
func (s *Store) ListBuildings() ([]Building, error) {
	names, err := s.base.List(buildingPrefix)
	if err != nil {
		return nil, err
	}
	buildings := make([]Building, 0, len(names))
	for _, name := range names {
		building, err := store.ReadJSON[Building](s.base, name)
		if err != nil {
			return nil, err
		}
		buildings = append(buildings, building)
	}
	return buildings, nil
}

// SaveZone persists a zone record.
func (s *Store) SaveZone(zone Zone) error {
	return store.WriteJSON(s.base, zonePrefix+zone.ID+".json", zone)
}

// ListZones returns zones belonging to a building.
func (s *Store) ListZones(buildingID string) ([]Zone, error) {
	names, err := s.base.List(zonePrefix)
	if err != nil {
		return nil, err
	}
	zones := make([]Zone, 0, len(names))
	for _, name := range names {
		zone, err := store.ReadJSON[Zone](s.base, name)
		if err != nil {
			return nil, err
		}
		if zone.BuildingID == buildingID {
			zones = append(zones, zone)
		}
	}
	return zones, nil
}
