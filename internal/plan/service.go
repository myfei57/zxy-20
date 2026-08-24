package plan

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/clock"
	"bms/internal/room"
)

// Service is the plan component entry point.
type Service struct {
	plans       *PlanStore
	switcher    *Switcher
	distributor *Distributor
	rooms       *room.Service
	clock       clock.Clock
}

// NewService wires the plan component over its stores and collaborators.
func NewService(plans *PlanStore, switcher *Switcher, distributor *Distributor, rooms *room.Service, clock clock.Clock) *Service {
	return &Service{plans: plans, switcher: switcher, distributor: distributor, rooms: rooms, clock: clock}
}

// Create registers a new draft schedule plan.
func (svc *Service) Create(name, buildingID string) (SchedulePlan, error) {
	if name == "" {
		return SchedulePlan{}, fmt.Errorf("plan name is required")
	}
	p := SchedulePlan{
		ID:         uuid.NewString(),
		Name:       name,
		BuildingID: buildingID,
		Version:    1,
		CreatedAt:  svc.clock.Now(),
	}
	if err := svc.plans.Save(p); err != nil {
		return SchedulePlan{}, err
	}
	return p, nil
}

// Plans lists every schedule plan.
func (svc *Service) Plans() ([]SchedulePlan, error) {
	return svc.plans.List()
}

// Plan loads one schedule plan by id.
func (svc *Service) Plan(id string) (SchedulePlan, error) {
	return svc.plans.Get(id)
}

// Switch activates the next version of a plan.
func (svc *Service) Switch(planID string) error {
	p, err := svc.plans.Get(planID)
	if err != nil {
		return err
	}
	p.Version++
	return svc.switcher.Switch(p)
}

// Distribute delivers a plan to every room of its building and returns the
// number of rooms that received it.
func (svc *Service) Distribute(planID, buildingID string) (int, error) {
	rooms, err := svc.rooms.RoomsByBuilding(buildingID)
	if err != nil {
		return 0, err
	}
	ids := make([]string, 0, len(rooms))
	for _, r := range rooms {
		ids = append(ids, r.ID)
	}
	if err := svc.distributor.Distribute(planID, ids); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// Cursor returns the distribution cursor of a plan.
func (svc *Service) Cursor(planID string) (int, error) {
	return svc.plans.Cursor(planID)
}

// ActiveState returns the durable activation marker of a plan.
func (svc *Service) ActiveState(planID string) (int, bool, error) {
	return svc.plans.Active(planID)
}

// Acked reports whether a room acknowledged the plan.
func (svc *Service) Acked(roomID, planID string) bool {
	return svc.distributor.Acked(roomID, planID)
}
