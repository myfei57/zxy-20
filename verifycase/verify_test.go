package verifycase

import (
	"testing"

	"bms/internal/clock"
	"bms/internal/device"
)

// TestDeviceReplaySkipsAckedCommands verifies a reconnecting device never
// receives a command whose acknowledgement is already recorded.
func TestDeviceReplaySkipsAckedCommands(t *testing.T) {
	inner := newMemStore()
	clk := clock.SystemClock{}

	svc := device.NewService(device.NewDeviceStore(inner), device.NewCommandStore(inner), device.NewStateStore(inner), clk)
	dev, err := svc.Register("room-1", "三号新风机组", "ahu")
	if err != nil {
		t.Fatal(err)
	}
	cmd, err := svc.Send(dev.ID, "start")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Ack(dev.ID, cmd.ID); err != nil {
		t.Fatal(err)
	}
	replay, err := svc.Replay(dev.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(replay) != 0 {
		t.Fatalf("replay re-sent %d acknowledged command(s)", len(replay))
	}
}
