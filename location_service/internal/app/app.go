package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"location_service/internal/httpadapter"
	openapi "location_service/internal/httpadapter/generate"
	"location_service/internal/repo/locationrepo"
	"location_service/internal/service"
	"moul.io/chizap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type app struct {
	config          *Config
	locationService *service.LocationService
	server          *http.Server
	logger          *zap.Logger
}

func New(config *Config, logger *zap.Logger) (App, error) {
	pgxPool, err := initDB(context.Background(), &config.Db)
	if err != nil {
		return nil, err
	}

	locationRepo := locationrepo.New(pgxPool, logger)
	locationService := service.New(locationRepo, logger)

	return &app{
		config:          config,
		locationService: locationService,
		logger:          logger,
	}, nil
}

func initDB(ctx context.Context, config *DatabaseConfig) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// migrations

	m, err := migrate.New(config.MigrationsDir, config.DSN)
	if err != nil {
		return nil, err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	if err := m.Up(); err != nil {
		return nil, err
	}

	return pool, nil
}

func (a *app) newHttpServer() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(chizap.New(a.logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	getDriversCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "location", Name: "get_drivers_request", Help: "Increment for /drivers endpoint",
	})

	updateDriverLocationCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "location", Name: "update_driver_request", Help: "Increment for /drivers/{driver_id}/location endpoint",
	})

	locationServer := httpadapter.New(*a.locationService, a.logger, getDriversCounter, updateDriverLocationCounter)

	petStoreStrictHandler := openapi.NewStrictHandler(locationServer, nil)
	openapi.HandlerFromMuxWithBaseURL(petStoreStrictHandler, router, a.config.App.BasePath)

	a.server = &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(a.config.App.Port)),
	}
}

func (a *app) Shutdown(ctx context.Context) error {
	<-ctx.Done()

	done := make(chan bool)
	a.logger.Info("Server is shutting down...")

	go func() {
		if err := a.server.Shutdown(context.Background()); err != nil {
			a.logger.Error("Could not gracefully shutdown the userHandler: ", zap.Error(err))
		}

		a.logger.Info("Server stopped ")
		close(done)
	}()

	<-done
	return nil
}

func (a *app) Serve() error {
	a.newHttpServer()

	// metrics for prometheus
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9000", nil) //nolint:errcheck

	done := make(chan os.Signal, 1)
	a.logger.Info(
		"Server is starting now...",
		zap.String("Port", strconv.Itoa(a.config.App.Port)),
	)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Could not listen on port: "+strconv.Itoa(a.config.App.Port), zap.Error(err))
		}
	}()

	<-done

	return nil
}
