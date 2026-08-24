package alert

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/clock"
)

// Service is the alert component entry point.
type Service struct {
	rules   *RuleStore
	store   *EventStore
	checker *Checker
}

// NewService wires the alert component over its stores.
func NewService(rules *RuleStore, store *EventStore, clock clock.Clock) *Service {
	return &Service{rules: rules, store: store, checker: NewChecker(rules, store, clock)}
}

// CreateRule registers a new threshold rule.
func (svc *Service) CreateRule(roomID, kind string, value float64) (ThresholdRule, error) {
	if value <= 0 {
		return ThresholdRule{}, fmt.Errorf("threshold must be positive")
	}
	rule := ThresholdRule{ID: uuid.NewString(), RoomID: roomID, Kind: kind, Value: value, Version: 1}
	if err := svc.rules.Save(rule); err != nil {
		return ThresholdRule{}, err
	}
	return rule, nil
}

// Rules lists every threshold rule.
func (svc *Service) Rules() ([]ThresholdRule, error) {
	return svc.rules.List()
}

// Activate marks a rule version active so it applies to new readings.
func (svc *Service) Activate(ruleID string) (ThresholdRule, error) {
	rules, err := svc.rules.List()
	if err != nil {
		return ThresholdRule{}, err
	}
	rule, err := RuleByID(rules, ruleID)
	if err != nil {
		return ThresholdRule{}, err
	}
	rule.Active = true
	rule.Version++
	if err := svc.rules.Save(rule); err != nil {
		return ThresholdRule{}, err
	}
	return rule, nil
}

// Events returns every triggered alert.
func (svc *Service) Events() ([]AlertEvent, error) {
	return svc.store.List()
}

// Checker exposes alert evaluation to the meter component.
func (svc *Service) Checker() *Checker {
	return svc.checker
}
