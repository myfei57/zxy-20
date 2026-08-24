package device

import (
	"fmt"

	"github.com/google/uuid"

	"bms/internal/clock"
)

// Service is the device component entry point.
type Service struct {
	devices *DeviceStore
	records *CommandStore
	state   *StateStore
	sender  *Sender
	replay  *Replayer
	clock   clock.Clock
}

// NewService wires the device component over its stores.
func NewService(devices *DeviceStore, records *CommandStore, state *StateStore, clock clock.Clock) *Service {
	return &Service{
		devices: devices,
		records: records,
		state:   state,
		sender:  NewSender(records, state, clock),
		replay:  NewReplayer(records),
		clock:   clock,
	}
}

// Register adds a new device to a room.
func (svc *Service) Register(roomID, name, kind string) (Device, error) {
	if name == "" {
		return Device{}, fmt.Errorf("device name is required")
	}
	now := svc.clock.Now()
	dev := Device{
		ID:        uuid.NewString(),
		RoomID:    roomID,
		Name:      name,
		Kind:      kind,
		State:     StateOff,
		UpdatedAt: now,
	}
	if err := svc.devices.Save(dev); err != nil {
		return Device{}, err
	}
	if err := svc.state.Init(dev.ID, StateOff, now); err != nil {
		return Device{}, err
	}
	return dev, nil
}

// Devices lists every registered device.
func (svc *Service) Devices() ([]Device, error) {
	return svc.devices.List()
}

// Send issues a command to a device.
func (svc *Service) Send(deviceID, command string) (*CommandRecord, error) {
	return svc.sender.Send(deviceID, command)
}

// Ack records a device acknowledgement and advances its state.
func (svc *Service) Ack(deviceID, commandID string) error {
	rec, err := svc.records.Ack(deviceID, commandID, svc.clock.Now())
	if err != nil {
		return err
	}
	return svc.state.SetState(deviceID, stateForCommand(rec.Command), svc.clock.Now())
}

// Replay returns commands to re-send after a device reconnect.
func (svc *Service) Replay(deviceID string) ([]CommandRecord, error) {
	return svc.replay.Replay(deviceID)
}

// CurrentState returns the live device state.
func (svc *Service) CurrentState(deviceID string) string {
	return svc.state.CurrentState(deviceID)
}

// Commands returns the durable command trace of a device.
func (svc *Service) Commands(deviceID string) ([]CommandRecord, error) {
	return svc.records.ListByDevice(deviceID)
}

func stateForCommand(command string) string {
	switch command {
	case "start":
		return StateRunning
	case "stop":
		return StateStopped
	case "fault":
		return StateFault
	default:
		return StateStarting
	}
}
