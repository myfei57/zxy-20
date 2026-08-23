package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) routes() {
	s.router.Use(withLogging, withRecover)
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/rooms", http.StatusFound)
	})
	s.router.Get("/rooms", s.page("rooms"))
	s.router.Get("/plans", s.page("plans"))
	s.router.Get("/meters", s.page("meters"))
	s.router.Get("/audit", s.page("audit"))

	s.router.Route("/api", func(r chi.Router) {
		r.Get("/buildings", s.listBuildings)
		r.Post("/buildings", s.createBuilding)
		r.Get("/zones", s.listZones)
		r.Post("/zones", s.createZone)

		r.Get("/rooms", s.listRooms)
		r.Post("/rooms", s.createRoom)
		r.Get("/rooms/{id}", s.viewRoom)
		r.Post("/rooms/{id}/setpoint", s.setRoomSetpoint)
		r.Post("/rooms/{id}/bind", s.bindRoom)
		r.Post("/rooms/{id}/device", s.attachRoomDevice)

		r.Get("/devices", s.listDevices)
		r.Post("/devices", s.createDevice)
		r.Get("/devices/{id}/commands", s.deviceCommands)
		r.Post("/devices/{id}/commands", s.sendCommand)
		r.Post("/devices/{id}/ack", s.ackCommand)
		r.Post("/devices/{id}/replay", s.replayDevice)

		r.Get("/plans", s.listPlans)
		r.Post("/plans", s.createPlan)
		r.Post("/plans/{id}/switch", s.switchPlan)
		r.Post("/plans/{id}/distribute", s.distributePlan)
		r.Get("/plans/{id}/ack/{roomID}", s.planAcked)

		r.Get("/meters", s.listMeters)
		r.Post("/meters", s.createMeter)
		r.Post("/meters/{id}/read", s.collectReading)
		r.Post("/meters/{id}/ingest", s.ingestReading)
		r.Get("/meters/{id}/readings", s.meterReadings)
		r.Get("/meters/{id}/windows", s.meterWindows)
		r.Post("/meters/windows/{id}/close", s.closeWindow)
		r.Get("/meters/windows/{id}", s.windowState)

		r.Get("/alerts", s.listAlerts)
		r.Get("/alerts/rules", s.listAlertRules)
		r.Post("/alerts/rules", s.createAlertRule)
		r.Post("/alerts/rules/{id}/activate", s.activateAlertRule)

		r.Get("/quotas", s.listQuotas)
		r.Post("/quotas", s.createQuota)
		r.Put("/quotas/{id}/limit", s.setQuotaLimit)

		r.Get("/audit", s.listAudit)
		r.Get("/reports/daily", s.dailyReport)
	})
}

func (s *Server) page(name string) http.HandlerFunc {
	body, ok := pages[name]
	if !ok {
		return func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(body))
	}
}
