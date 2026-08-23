package console

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"bms/internal/alert"
	"bms/internal/audit"
	"bms/internal/config"
	"bms/internal/device"
	"bms/internal/meter"
	"bms/internal/ns"
	"bms/internal/plan"
	"bms/internal/quota"
	"bms/internal/report"
	"bms/internal/room"
)

// Deps bundles every component the console talks to.
type Deps struct {
	Buildings *ns.Service
	Rooms     *room.Service
	Plans     *plan.Service
	Devices   *device.Service
	Meters    *meter.Service
	Alerts    *alert.Service
	Quotas    *quota.Service
	Audit     *audit.Service
	Reports   *report.Builder
}

// Server is the HTTP control console.
type Server struct {
	cfg    config.Config
	deps   Deps
	router *chi.Mux
	http   *http.Server
}

// NewServer builds the console router and HTTP server.
func NewServer(cfg config.Config, deps Deps) *Server {
	server := &Server{cfg: cfg, deps: deps, router: chi.NewRouter()}
	server.routes()
	return server
}

// Start begins serving HTTP requests.
func (s *Server) Start() error {
	s.http = &http.Server{Addr: s.cfg.Addr, Handler: s.router}
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.http == nil {
		return nil
	}
	return s.http.Shutdown(ctx)
}
