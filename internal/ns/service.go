package ns

import (
	"fmt"

	"github.com/google/uuid"
)

// Service is the namespace entry point used by the console and seed.
type Service struct {
	store *Store
}

// NewService wires a namespace service over the given store.
func NewService(store *Store) *Service {
	return &Service{store: store}
}

// CreateBuilding registers a new building.
func (svc *Service) CreateBuilding(name, address string) (Building, error) {
	if name == "" {
		return Building{}, fmt.Errorf("building name is required")
	}
	building := Building{ID: uuid.NewString(), Name: name, Address: address}
	if err := svc.store.SaveBuilding(building); err != nil {
		return Building{}, err
	}
	return building, nil
}

// CreateZone adds a zone to an existing building.
func (svc *Service) CreateZone(buildingID, name, floor string) (Zone, error) {
	if _, err := svc.store.GetBuilding(buildingID); err != nil {
		return Zone{}, fmt.Errorf("building %s: %w", buildingID, err)
	}
	if name == "" {
		return Zone{}, fmt.Errorf("zone name is required")
	}
	zone := Zone{ID: uuid.NewString(), BuildingID: buildingID, Name: name, Floor: floor}
	if err := svc.store.SaveZone(zone); err != nil {
		return Zone{}, err
	}
	return zone, nil
}

// Buildings lists all buildings.
func (svc *Service) Buildings() ([]Building, error) {
	return svc.store.ListBuildings()
}

// Zones lists zones for a building.
func (svc *Service) Zones(buildingID string) ([]Zone, error) {
	return svc.store.ListZones(buildingID)
}
