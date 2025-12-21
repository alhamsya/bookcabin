package rest

import (
	"context"

	"github.com/alhamsya/bookcabin/internal/adapter/airline/airasia"
	"github.com/alhamsya/bookcabin/internal/adapter/airline/batik"
	"github.com/alhamsya/bookcabin/internal/adapter/airline/garuda"
	"github.com/alhamsya/bookcabin/internal/adapter/airline/lion"

	"github.com/alhamsya/bookcabin/internal/adapter/handler/rest"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alhamsya/bookcabin/internal/adapter/redis"
	"github.com/alhamsya/bookcabin/internal/core/service/flight"
	"github.com/alhamsya/bookcabin/lib/manager/config"
	"github.com/alhamsya/bookcabin/lib/manager/protocol"

	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS to match Linux container CPU quota.
)

func RunApp(ctx context.Context) error { //nolint:nolintlint,funlen
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	/* === GENERAL === */
	cfg := config.GetConfig(ctx)

	/* === HANDLER === */
	/* === DATABASE === */

	// Init cache service
	cache, err := redis.New(ctx, cfg)
	if err != nil {
		log.Fatalln("[Rest] redis not running", err)
	}
	defer cache.Close()

	// init Airline service
	airAsiaRepo := airasia.New(&airasia.Airline{
		CfgUpstream: cfg.Static.Upstream,
	})
	batikRepo := batik.New(&batik.Airline{
		CfgUpstream: cfg.Static.Upstream,
	})
	lionRepo := lion.New(&lion.Airline{
		CfgUpstream: cfg.Static.Upstream,
	})
	garudaRepo := garuda.New(&garuda.Airline{
		CfgUpstream: cfg.Static.Upstream,
	})

	// Dependency injection
	// Flight
	flightService := flight.NewFlightService(&flight.Service{
		Cfg:         cfg,
		Cache:       cache,
		AirAsiaRepo: airAsiaRepo,
		BatikRepo:   batikRepo,
		GarudaRepo:  garudaRepo,
		LionRepo:    lionRepo,
	})

	// Init router
	server := protocol.Rest(ctx, &protocol.RESTService{
		Cfg: cfg,
		Interactor: &rest.Interactor{
			FlightService: flightService,
		},
	})
	if err = server.Run(); err != nil {
		log.Fatalln("[Rest] service not running", err)
	}

	return nil
}
