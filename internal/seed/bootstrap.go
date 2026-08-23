package seed

import (
	"fmt"

	"bms/internal/alert"
	"bms/internal/device"
	"bms/internal/meter"
	"bms/internal/ns"
	"bms/internal/plan"
	"bms/internal/quota"
	"bms/internal/report"
	"bms/internal/room"
)

// Services bundles every component needed to bootstrap demo data.
type Services struct {
	NS      *ns.Service
	Rooms   *room.Service
	Devices *device.Service
	Plans   *plan.Service
	Meters  *meter.Service
	Alerts  *alert.Service
	Quotas  *quota.Service
	Reports *report.Builder
}

// Bootstrap writes demo data once when the store is empty.
func Bootstrap(svc Services) error {
	buildings, err := svc.NS.Buildings()
	if err != nil {
		return err
	}
	if len(buildings) > 0 {
		return nil
	}
	building, err := svc.NS.CreateBuilding("A栋", "朝阳区科技园 1 号")
	if err != nil {
		return fmt.Errorf("create demo building: %w", err)
	}
	zoneOffice, err := svc.NS.CreateZone(building.ID, "三层办公区", "3F")
	if err != nil {
		return err
	}
	zoneMeeting, err := svc.NS.CreateZone(building.ID, "五层会议区", "5F")
	if err != nil {
		return err
	}
	room1, err := svc.Rooms.Create(building.ID, zoneOffice.ID, "三层会客室", room.ModeAuto, "")
	if err != nil {
		return err
	}
	room2, err := svc.Rooms.Create(building.ID, zoneOffice.ID, "三层开放办公区", room.ModeAuto, "")
	if err != nil {
		return err
	}
	room3, err := svc.Rooms.Create(building.ID, zoneMeeting.ID, "五层大会议室", room.ModeCool, "")
	if err != nil {
		return err
	}
	device1, err := svc.Devices.Register(room1.ID, "三层风机盘管 01", "fcu")
	if err != nil {
		return err
	}
	device2, err := svc.Devices.Register(room2.ID, "三层新风机组 01", "ahu")
	if err != nil {
		return err
	}
	device3, err := svc.Devices.Register(room3.ID, "五层冷机阀组", "vav")
	if err != nil {
		return err
	}
	for _, binding := range []struct{ roomID, deviceID string }{
		{room1.ID, device1.ID},
		{room2.ID, device2.ID},
		{room3.ID, device3.ID},
	} {
		if err := svc.Rooms.AttachDevice(binding.roomID, binding.deviceID); err != nil {
			return err
		}
	}
	meter1, err := svc.Meters.CreateMeter(room1.ID, "三层电表 01")
	if err != nil {
		return err
	}
	meter2, err := svc.Meters.CreateMeter(room2.ID, "三层电表 02")
	if err != nil {
		return err
	}
	meter3, err := svc.Meters.CreateMeter(room3.ID, "五层电表 01")
	if err != nil {
		return err
	}
	for _, item := range []struct {
		roomID string
		limit  float64
	}{
		{room1.ID, 500},
		{room2.ID, 800},
		{room3.ID, 300},
	} {
		if _, err := svc.Quotas.Create(item.roomID, item.limit); err != nil {
			return err
		}
	}
	rule, err := svc.Alerts.CreateRule(room1.ID, "power", 60)
	if err != nil {
		return err
	}
	if _, err := svc.Alerts.Activate(rule.ID); err != nil {
		return err
	}
	demoPlan, err := svc.Plans.Create("夏季工作日计划", building.ID)
	if err != nil {
		return err
	}
	if err := svc.Plans.Switch(demoPlan.ID); err != nil {
		return err
	}
	if _, err := svc.Plans.Distribute(demoPlan.ID, building.ID); err != nil {
		return err
	}
	if _, err := svc.Meters.Collect(meter1.ID, 12.3); err != nil {
		return err
	}
	if err := svc.Meters.Ingest(meter2.ID, 8.6); err != nil {
		return err
	}
	if err := svc.Meters.Ingest(meter3.ID, 21.4); err != nil {
		return err
	}
	window, err := svc.Meters.OpenWindow(meter1.ID)
	if err != nil {
		return err
	}
	if err := svc.Meters.CloseAfterSummary(window.ID); err != nil {
		return err
	}
	command, err := svc.Devices.Send(device2.ID, "start")
	if err != nil {
		return err
	}
	if err := svc.Devices.Ack(device2.ID, command.ID); err != nil {
		return err
	}
	if _, err := svc.Reports.Build(today()); err != nil {
		return err
	}
	return nil
}
