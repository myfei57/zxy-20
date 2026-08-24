package alert

import (
	"fmt"

	"bms/internal/store"
)

const (
	rulePrefix  = "alert/rules/"
	eventPrefix = "alert/events/"
)

// RuleStore persists threshold rules durably.
type RuleStore struct {
	base store.Store
}

// NewRuleStore wraps a durable store for threshold rules.
func NewRuleStore(base store.Store) *RuleStore {
	return &RuleStore{base: base}
}

// Save persists one threshold rule.
func (s *RuleStore) Save(rule ThresholdRule) error {
	return store.WriteJSON(s.base, rulePrefix+rule.ID+".json", rule)
}

// List returns every threshold rule sorted by record name.
func (s *RuleStore) List() ([]ThresholdRule, error) {
	names, err := s.base.List(rulePrefix)
	if err != nil {
		return nil, err
	}
	rules := make([]ThresholdRule, 0, len(names))
	for _, name := range names {
		rule, err := store.ReadJSON[ThresholdRule](s.base, name)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// ListActive returns the rules that are currently activated.
func (s *RuleStore) ListActive() ([]ThresholdRule, error) {
	rules, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]ThresholdRule, 0, len(rules))
	for _, rule := range rules {
		if rule.Active {
			out = append(out, rule)
		}
	}
	return out, nil
}

// EventStore persists triggered alerts durably.
type EventStore struct {
	base store.Store
}

// NewEventStore wraps a durable store for alert events.
func NewEventStore(base store.Store) *EventStore {
	return &EventStore{base: base}
}

// Append persists one alert event.
func (s *EventStore) Append(event AlertEvent) error {
	return store.WriteJSON(s.base, eventPrefix+event.ID+".json", event)
}

// List returns every alert event sorted by record name.
func (s *EventStore) List() ([]AlertEvent, error) {
	names, err := s.base.List(eventPrefix)
	if err != nil {
		return nil, err
	}
	events := make([]AlertEvent, 0, len(names))
	for _, name := range names {
		event, err := store.ReadJSON[AlertEvent](s.base, name)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// RuleByID finds a rule by id, used by activation flows.
func RuleByID(rules []ThresholdRule, id string) (ThresholdRule, error) {
	for _, rule := range rules {
		if rule.ID == id {
			return rule, nil
		}
	}
	return ThresholdRule{}, fmt.Errorf("rule %s not found", id)
}
