package device

import (
	"errors"
	"testing"
	"time"

	"bms/internal/store"
)

type fakeClock struct{ t time.Time }

func (c fakeClock) Now() time.Time { return c.t }

// commandsFailStore lets device-state writes through but fails every write
// under the command-record prefix, isolating the record-persistence failure.
type commandsFailStore struct {
	store.Store
	err error
}

func (s *commandsFailStore) Write(name string, data []byte) error {
	if len(name) >= len(commandPrefix) && name[:len(commandPrefix)] == commandPrefix {
		return s.err
	}
	return s.Store.Write(name, data)
}

// TestSendWritesRecordBeforeState asserts the console invariant: whenever a
// device shows "sent" the matching command record has already been durably
// written. A failed record write must fail Send without touching device state.
func TestSendWritesRecordBeforeState(t *testing.T) {
	setup := func() (*Sender, *StateStore, *CommandStore) {
		base, err := store.NewFileStore(t.TempDir())
		if err != nil {
			t.Fatalf("open store: %v", err)
		}
		if err := (&DeviceStore{base: base}).Save(Device{
			ID: "chiller-2", RoomID: "plant", Name: "二号冷水机", Kind: "chiller", State: StateOff,
		}); err != nil {
			t.Fatalf("seed device: %v", err)
		}
		state := NewStateStore(base)
		if err := state.Init("chiller-2", StateOff, time.Now()); err != nil {
			t.Fatalf("init state: %v", err)
		}
		records := NewCommandStore(base)
		return NewSender(records, state, fakeClock{t: time.Unix(1700000000, 0)}), state, records
	}

	t.Run("happy path persists record before marking sent", func(t *testing.T) {
		sender, _, records := setup()
		rec, err := sender.Send("chiller-2", "stop")
		if err != nil {
			t.Fatalf("send: %v", err)
		}
		got, err := records.Get(rec.ID)
		if err != nil {
			t.Fatalf("record missing after send: %v", err)
		}
		if got.Status != StatusSent || got.Command != "stop" {
			t.Fatalf("record = %+v, want status=%q command=%q", got, StatusSent, "stop")
		}
	})

	t.Run("record write failure leaves device state untouched", func(t *testing.T) {
		base, err := store.NewFileStore(t.TempDir())
		if err != nil {
			t.Fatalf("open store: %v", err)
		}
		if err := (&DeviceStore{base: base}).Save(Device{
			ID: "chiller-2", RoomID: "plant", Name: "二号冷水机", Kind: "chiller", State: StateOff,
		}); err != nil {
			t.Fatalf("seed device: %v", err)
		}
		state := NewStateStore(base)
		if err := state.Init("chiller-2", StateOff, time.Now()); err != nil {
			t.Fatalf("init state: %v", err)
		}

		boom := errors.New("disk full")
		flaky := &commandsFailStore{Store: base, err: boom}
		records := NewCommandStore(flaky)
		sender := NewSender(records, state, fakeClock{t: time.Unix(1700000001, 0)})

		beforeState := state.CurrentState("chiller-2")
		// Capture LastCommandID before the failed attempt.
		before, _ := store.ReadJSON[Device](base, statePrefix+"chiller-2.json")

		_, err = sender.Send("chiller-2", "start")
		if err == nil {
			t.Fatal("expected send to fail when record write fails")
		}

		after, _ := store.ReadJSON[Device](base, statePrefix+"chiller-2.json")
		if after.LastCommandID != before.LastCommandID {
			t.Fatalf("device LastCommandID changed despite record write failure: before=%q after=%q",
				before.LastCommandID, after.LastCommandID)
		}
		if state.CurrentState("chiller-2") != beforeState {
			t.Fatalf("device state changed despite record write failure")
		}

		list, err := records.ListByDevice("chiller-2")
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		for _, r := range list {
			if r.Command == "start" && r.SentAt.Unix() == 1700000001 {
				t.Fatalf("a 'start' command from the failed send persisted: %+v", r)
			}
		}
	})
}
