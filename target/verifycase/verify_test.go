package verifycase

import (
	"testing"

	"bms/internal/clock"
	"bms/internal/room"
)

// TestRoomSetpointAfterStoreDurable verifies the displayed setpoint only
// changes after the durable room store accepts the write.
func TestRoomSetpointAfterStoreDurable(t *testing.T) {
	inner := newMemStore()
	clk := clock.SystemClock{}

	roomStore := room.NewStore(inner)
	r := room.Room{ID: "room-1", BuildingID: "b1", Name: "五层会议室", Setpoint: 26}
	if err := roomStore.SaveRoom(r); err != nil {
		t.Fatal(err)
	}
	cache := room.NewCache()
	cache.Set(r)

	failing := &failStore{inner: inner, marker: "rooms/"}
	svc := room.NewService(room.NewStore(failing), cache, nil, clk)
	if err := svc.SetSetpoint("room-1", 24); err == nil {
		t.Fatal("setpoint update must fail when the store write fails")
	}
	got, ok := cache.Get("room-1")
	if !ok {
		t.Fatal("room missing from cache")
	}
	if got.Setpoint != 26 {
		t.Fatalf("displayed setpoint %v changed although the durable write failed", got.Setpoint)
	}
}
