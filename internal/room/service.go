package room

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/clock"
	"bms/internal/device"
)

// Service is the room component entry point.
type Service struct {
	store     *Store
	cache     *Cache
	setpoints *SetpointService
	binder    *Binder
	query     *Query
	clock     clock.Clock
}

// NewService wires the room component over its stores and the device service.
func NewService(store *Store, cache *Cache, devices *device.Service, clock clock.Clock) *Service {
	return &Service{
		store:     store,
		cache:     cache,
		setpoints: NewSetpointService(store, cache, clock),
		binder:    NewBinder(store, cache, clock),
		query:     NewQuery(store, cache, devices),
		clock:     clock,
	}
}

// Binder exposes the room binder to the plan switcher.
func (svc *Service) Binder() *Binder {
	return svc.binder
}

// Create registers a new room.
func (svc *Service) Create(buildingID, zoneID, name string, mode Mode, deviceID string) (Room, error) {
	if name == "" {
		return Room{}, fmt.Errorf("room name is required")
	}
	r := Room{
		ID:         uuid.NewString(),
		BuildingID: buildingID,
		ZoneID:     zoneID,
		Name:       name,
		Setpoint:   26,
		Mode:       mode,
		DeviceID:   deviceID,
		UpdatedAt:  svc.clock.Now(),
	}
	if err := svc.store.SaveRoom(r); err != nil {
		return Room{}, err
	}
	svc.cache.Set(r)
	return r, nil
}

// Rooms returns live views for every room.
func (svc *Service) Rooms() ([]View, error) {
	return svc.query.All()
}

// View returns the live view of one room.
func (svc *Service) View(roomID string) (View, error) {
	return svc.query.ByID(roomID)
}

// SetSetpoint updates a room setpoint durably before showing the new value.
func (svc *Service) SetSetpoint(roomID string, value float64) error {
	return svc.setpoints.Set(roomID, value)
}

// BindToPlan binds one room to a schedule plan.
func (svc *Service) BindToPlan(roomID, planID string) error {
	return svc.binder.BindToPlan(roomID, planID)
}

// RoomsByBuilding returns rooms inside a building.
func (svc *Service) RoomsByBuilding(buildingID string) ([]Room, error) {
	return svc.store.ListByBuilding(buildingID)
}

// AttachDevice binds an installed device to a room durably.
func (svc *Service) AttachDevice(roomID, deviceID string) error {
	r, ok := svc.cache.Get(roomID)
	if !ok {
		var err error
		r, err = svc.store.GetRoom(roomID)
		if err != nil {
			return fmt.Errorf("attach device to %s: %w", roomID, err)
		}
	}
	updated := r
	updated.DeviceID = deviceID
	updated.UpdatedAt = svc.clock.Now()
	if err := svc.store.SaveRoom(updated); err != nil {
		return fmt.Errorf("persist device binding %s: %w", roomID, err)
	}
	svc.cache.Set(updated)
	return nil
}
