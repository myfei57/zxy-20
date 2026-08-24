package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bms/internal/alert"
	"bms/internal/audit"
	"bms/internal/clock"
	"bms/internal/config"
	"bms/internal/console"
	"bms/internal/device"
	"bms/internal/meter"
	"bms/internal/ns"
	"bms/internal/plan"
	"bms/internal/quota"
	"bms/internal/report"
	"bms/internal/room"
	"bms/internal/seed"
	"bms/internal/store"
)

func main() {
	cfg := config.Load()
	var baseStore *store.FileStore
	var err error
	baseStore, err = store.NewFileStore(cfg.DataDir)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	var clk clock.SystemClock = clock.NewSystemClock()

	nsService := ns.NewService(ns.NewStore(baseStore))

	roomStore := room.NewStore(baseStore)
	roomCache := room.NewCache()

	deviceService := device.NewService(
		device.NewDeviceStore(baseStore),
		device.NewCommandStore(baseStore),
		device.NewStateStore(baseStore),
		clk,
	)
	roomService := room.NewService(roomStore, roomCache, deviceService, clk)

	auditService := audit.NewService(audit.NewStore(baseStore), clk)
	planStore := plan.NewStore(baseStore)
	switcher := plan.NewSwitcher(planStore, roomService.Binder(), auditService.Recorder(), clk)
	distributor := plan.NewDistributor(planStore, meter.NewAckStore(baseStore, clk))
	planService := plan.NewService(planStore, switcher, distributor, roomService, clk)

	quotaService := quota.NewService(quota.NewStore(baseStore))
	alertService := alert.NewService(alert.NewRuleStore(baseStore), alert.NewEventStore(baseStore), clk)

	readingStore := meter.NewReadingStore(baseStore)
	meterService := meter.NewService(
		meter.NewMeterStore(baseStore),
		readingStore,
		meter.NewCursorStore(baseStore, readingStore),
		meter.NewWindowStore(baseStore, clk),
		meter.NewSummaryStore(baseStore),
		quotaService.Checker(),
		alertService.Checker(),
		clk,
	)
	reportService := report.NewBuilder(meterService)

	if cfg.SeedDemo {
		if err := seed.Bootstrap(seed.Services{
			NS:      nsService,
			Rooms:   roomService,
			Devices: deviceService,
			Plans:   planService,
			Meters:  meterService,
			Alerts:  alertService,
			Quotas:  quotaService,
			Reports: reportService,
		}); err != nil {
			log.Fatalf("seed demo data: %v", err)
		}
	}

	server := console.NewServer(cfg, console.Deps{
		Buildings: nsService,
		Rooms:     roomService,
		Plans:     planService,
		Devices:   deviceService,
		Meters:    meterService,
		Alerts:    alertService,
		Quotas:    quotaService,
		Audit:     auditService,
		Reports:   reportService,
	})

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve: %v", err)
		}
	}()
	log.Printf("BMS console listening on %s (data: %s)", cfg.Addr, baseStore.Root())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
