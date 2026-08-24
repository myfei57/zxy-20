package verifycase

import (
	"testing"

	"bms/internal/clock"
	"bms/internal/device"
	"bms/internal/room"
)

// TestRoomQueryUsesCurrentDeviceState verifies the room page reads the live
// device state instead of a snapshot captured before the command ack.
func TestRoomQueryUsesCurrentDeviceState(t *testing.T) {
	inner := newMemStore()
	clk := clock.SystemClock{}

	deviceSvc := device.NewService(device.NewDeviceStore(inner), device.NewCommandStore(inner), device.NewStateStore(inner), clk)
	roomStore := room.NewStore(inner)
	cache := room.NewCache()
	roomSvc := room.NewService(roomStore, cache, deviceSvc, clk)

	dev, err := deviceSvc.Register("room-1", "三号新风机组", "ahu")
	if err != nil {
		t.Fatal(err)
	}
	r := room.Room{ID: "room-1", BuildingID: "b1", Name: "三层开放办公区", DeviceID: dev.ID}
	if err := roomStore.SaveRoom(r); err != nil {
		t.Fatal(err)
	}
	cache.Set(r)

	start, err := deviceSvc.Send(dev.ID, "start")
	if err != nil {
		t.Fatal(err)
	}
	if err := deviceSvc.Ack(dev.ID, start.ID); err != nil {
		t.Fatal(err)
	}
	first, err := roomSvc.View("room-1")
	if err != nil {
		t.Fatal(err)
	}
	if first.DeviceState != device.StateRunning {
		t.Fatalf("expected running before stop, got %q", first.DeviceState)
	}

	stop, err := deviceSvc.Send(dev.ID, "stop")
	if err != nil {
		t.Fatal(err)
	}
	if err := deviceSvc.Ack(dev.ID, stop.ID); err != nil {
		t.Fatal(err)
	}
	second, err := roomSvc.View("room-1")
	if err != nil {
		t.Fatal(err)
	}
	if second.DeviceState != device.StateStopped {
		t.Fatalf("room page shows stale device state %q", second.DeviceState)
	}
}
