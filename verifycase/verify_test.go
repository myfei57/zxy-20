package verifycase

import (
	"testing"

	"bms/internal/clock"
	"bms/internal/device"
	"bms/internal/store"
)

// TestDeviceSentAfterCommandDurable verifies a device never shows a sent
// command while its command record is missing from the durable trace.
func TestDeviceSentAfterCommandDurable(t *testing.T) {
	inner := newMemStore()
	failing := &failStore{inner: inner, marker: "device/commands/"}
	clk := clock.SystemClock{}

	stateStore := device.NewStateStore(inner)
	svc := device.NewService(device.NewDeviceStore(inner), device.NewCommandStore(failing), stateStore, clk)
	dev, err := svc.Register("room-1", "二号冷水机", "chiller")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Send(dev.ID, "start"); err == nil {
		t.Fatal("send must fail when the command record cannot be written")
	}
	stored, err := store.ReadJSON[device.Device](inner, "device/states/"+dev.ID+".json")
	if err != nil {
		t.Fatal(err)
	}
	if stored.LastCommandID != "" {
		t.Fatalf("device shows sent command %s although no record is durable", stored.LastCommandID)
	}
}
