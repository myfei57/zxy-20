package console

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"bms/internal/audit"
	"bms/internal/room"
)

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, map[string]string{"error": message})
}

func readJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func fail(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, err.Error())
}

func (s *Server) listBuildings(w http.ResponseWriter, r *http.Request) {
	buildings, err := s.deps.Buildings.Buildings()
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildings)
}

func (s *Server) createBuilding(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	building, err := s.deps.Buildings.CreateBuilding(request.Name, request.Address)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, building)
}

func (s *Server) listZones(w http.ResponseWriter, r *http.Request) {
	buildingID := r.URL.Query().Get("building_id")
	zones, err := s.deps.Buildings.Zones(buildingID)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

func (s *Server) createZone(w http.ResponseWriter, r *http.Request) {
	var request struct {
		BuildingID string `json:"building_id"`
		Name       string `json:"name"`
		Floor      string `json:"floor"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	zone, err := s.deps.Buildings.CreateZone(request.BuildingID, request.Name, request.Floor)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, zone)
}

func (s *Server) listRooms(w http.ResponseWriter, r *http.Request) {
	views, err := s.deps.Rooms.Rooms()
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, views)
}

func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {
	var request struct {
		BuildingID string `json:"building_id"`
		ZoneID     string `json:"zone_id"`
		Name       string `json:"name"`
		Mode       string `json:"mode"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	created, err := s.deps.Rooms.Create(request.BuildingID, request.ZoneID, request.Name, room.Mode(request.Mode), "")
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) viewRoom(w http.ResponseWriter, r *http.Request) {
	view, err := s.deps.Rooms.View(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) setRoomSetpoint(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Value float64 `json:"value"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	if err := s.deps.Rooms.SetSetpoint(chi.URLParam(r, "id"), request.Value); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) bindRoom(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PlanID string `json:"plan_id"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	if err := s.deps.Rooms.BindToPlan(chi.URLParam(r, "id"), request.PlanID); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) attachRoomDevice(w http.ResponseWriter, r *http.Request) {
	var request struct {
		DeviceID string `json:"device_id"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	if err := s.deps.Rooms.AttachDevice(chi.URLParam(r, "id"), request.DeviceID); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) listDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.deps.Devices.Devices()
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, devices)
}

func (s *Server) createDevice(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RoomID string `json:"room_id"`
		Name   string `json:"name"`
		Kind   string `json:"kind"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	device, err := s.deps.Devices.Register(request.RoomID, request.Name, request.Kind)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, device)
}

func (s *Server) sendCommand(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Command string `json:"command"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	record, err := s.deps.Devices.Send(chi.URLParam(r, "id"), request.Command)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) ackCommand(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CommandID string `json:"command_id"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	if err := s.deps.Devices.Ack(chi.URLParam(r, "id"), request.CommandID); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) replayDevice(w http.ResponseWriter, r *http.Request) {
	records, err := s.deps.Devices.Replay(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) deviceCommands(w http.ResponseWriter, r *http.Request) {
	records, err := s.deps.Devices.Commands(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) listPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := s.deps.Plans.Plans()
	if err != nil {
		fail(w, err)
		return
	}
	for index := range plans {
		plan := &plans[index]
		version, active, err := s.deps.Plans.ActiveState(plan.ID)
		if err != nil {
			fail(w, err)
			return
		}
		cursor, err := s.deps.Plans.Cursor(plan.ID)
		if err != nil {
			fail(w, err)
			return
		}
		plan.Version = version
		plan.Active = active
		plan.Cursor = cursor
	}
	writeJSON(w, http.StatusOK, plans)
}

func (s *Server) createPlan(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name       string `json:"name"`
		BuildingID string `json:"building_id"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	plan, err := s.deps.Plans.Create(request.Name, request.BuildingID)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, plan)
}

func (s *Server) switchPlan(w http.ResponseWriter, r *http.Request) {
	if err := s.deps.Plans.Switch(chi.URLParam(r, "id")); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) distributePlan(w http.ResponseWriter, r *http.Request) {
	plan, err := s.deps.Plans.Plan(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	count, err := s.deps.Plans.Distribute(plan.ID, plan.BuildingID)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"delivered": count})
}

func (s *Server) planAcked(w http.ResponseWriter, r *http.Request) {
	acked := s.deps.Plans.Acked(chi.URLParam(r, "roomID"), chi.URLParam(r, "id"))
	writeJSON(w, http.StatusOK, map[string]bool{"acked": acked})
}

func (s *Server) listMeters(w http.ResponseWriter, r *http.Request) {
	meters, err := s.deps.Meters.Meters()
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, meters)
}

func (s *Server) createMeter(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RoomID string `json:"room_id"`
		Name   string `json:"name"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	meter, err := s.deps.Meters.CreateMeter(request.RoomID, request.Name)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, meter)
}

func (s *Server) collectReading(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Value float64 `json:"value"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	reading, err := s.deps.Meters.Collect(chi.URLParam(r, "id"), request.Value)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, reading)
}

func (s *Server) ingestReading(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Value float64 `json:"value"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	if err := s.deps.Meters.Ingest(chi.URLParam(r, "id"), request.Value); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) meterReadings(w http.ResponseWriter, r *http.Request) {
	readings, err := s.deps.Meters.Readings(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, readings)
}

func (s *Server) meterWindows(w http.ResponseWriter, r *http.Request) {
	windows, err := s.deps.Meters.Windows(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, windows)
}

func (s *Server) closeWindow(w http.ResponseWriter, r *http.Request) {
	if err := s.deps.Meters.CloseAfterSummary(chi.URLParam(r, "id")); err != nil {
		fail(w, err)
		return
	}
	closed, err := s.deps.Meters.WindowClosed(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "closed": closed})
}

func (s *Server) windowState(w http.ResponseWriter, r *http.Request) {
	closed, err := s.deps.Meters.WindowClosed(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"window_id": chi.URLParam(r, "id"), "closed": closed})
}

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	events, err := s.deps.Alerts.Events()
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.deps.Alerts.Rules()
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RoomID string  `json:"room_id"`
		Kind   string  `json:"kind"`
		Value  float64 `json:"value"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	rule, err := s.deps.Alerts.CreateRule(request.RoomID, request.Kind, request.Value)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, rule)
}

func (s *Server) activateAlertRule(w http.ResponseWriter, r *http.Request) {
	rule, err := s.deps.Alerts.Activate(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rule)
}

func (s *Server) listQuotas(w http.ResponseWriter, r *http.Request) {
	quotas, err := s.deps.Quotas.Quotas()
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, quotas)
}

func (s *Server) createQuota(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RoomID string  `json:"room_id"`
		Limit  float64 `json:"limit"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	quota, err := s.deps.Quotas.Create(request.RoomID, request.Limit)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, quota)
}

func (s *Server) setQuotaLimit(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Limit float64 `json:"limit"`
	}
	if err := readJSON(r, &request); err != nil {
		fail(w, err)
		return
	}
	quota, err := s.deps.Quotas.SetLimit(chi.URLParam(r, "id"), request.Limit)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, quota)
}

func (s *Server) listAudit(w http.ResponseWriter, r *http.Request) {
	kind := r.URL.Query().Get("kind")
	limitText := r.URL.Query().Get("limit")
	limit := 0
	if limitText != "" {
		if _, err := fmt.Sscanf(limitText, "%d", &limit); err != nil {
			fail(w, err)
			return
		}
	}
	var (
		events []audit.AuditEvent
		err    error
	)
	switch {
	case limit > 0:
		events, err = s.deps.Audit.Recent(limit)
	case kind != "":
		events, err = s.deps.Audit.EventsByKind(kind)
	default:
		events, err = s.deps.Audit.Events()
	}
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) dailyReport(w http.ResponseWriter, r *http.Request) {
	day := r.URL.Query().Get("day")
	if day == "" {
		day = time.Now().Format("2006-01-02")
	}
	reports, err := s.deps.Reports.Build(day)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, reports)
}
