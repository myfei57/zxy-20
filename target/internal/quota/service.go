package quota

import (
	"fmt"

	"github.com/google/uuid"
)

// Service is the quota component entry point.
type Service struct {
	store   *Store
	checker *Checker
}

// NewService wires the quota component over its store.
func NewService(store *Store) *Service {
	return &Service{store: store, checker: NewChecker(store)}
}

// Create assigns an energy quota to a room.
func (svc *Service) Create(roomID string, limit float64) (Quota, error) {
	if limit <= 0 {
		return Quota{}, fmt.Errorf("quota limit must be positive")
	}
	q := Quota{ID: uuid.NewString(), Scope: "room", RoomID: roomID, Limit: limit, Version: 1}
	if err := svc.store.Save(q); err != nil {
		return Quota{}, err
	}
	return q, nil
}

// Quotas lists every energy quota.
func (svc *Service) Quotas() ([]Quota, error) {
	return svc.store.List()
}

// SetLimit updates the durable limit of a quota and bumps its version.
func (svc *Service) SetLimit(quotaID string, limit float64) (Quota, error) {
	q, err := svc.byID(quotaID)
	if err != nil {
		return Quota{}, err
	}
	q.Limit = limit
	q.Version++
	if err := svc.store.Save(q); err != nil {
		return Quota{}, err
	}
	return q, nil
}

// Checker exposes the quota gate to the meter component.
func (svc *Service) Checker() *Checker {
	return svc.checker
}

func (svc *Service) byID(quotaID string) (Quota, error) {
	quotas, err := svc.store.List()
	if err != nil {
		return Quota{}, err
	}
	for _, q := range quotas {
		if q.ID == quotaID {
			return q, nil
		}
	}
	return Quota{}, fmt.Errorf("quota %s not found", quotaID)
}
