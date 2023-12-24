package app

import (
	"context"
	config "driver_service/configs"
	http3 "driver_service/configs/http"
	consumer2 "driver_service/configs/kafka/consumer"
	producer2 "driver_service/configs/kafka/producer"
	http2 "driver_service/internal/adapters/http"
	openapi "driver_service/internal/adapters/http/generate"
	trip_outbound2 "driver_service/internal/application/eventbus/producer/trip_outbound"
	"driver_service/internal/application/services/driver"
	"driver_service/internal/application/services/location"
	"driver_service/internal/application/services/trip"
	domainTrip "driver_service/internal/domain/repository/trip"
	"driver_service/internal/infrastracture/eventbus/consumer"
	kafka_consumers "driver_service/internal/infrastracture/eventbus/consumer/kafka"
	"driver_service/internal/infrastracture/eventbus/consumer/kafka/trip_outbound"
	trip_outbound3 "driver_service/internal/infrastracture/eventbus/producer/trip_outbound"
	"driver_service/internal/infrastracture/logger"
	infraTrip "driver_service/internal/infrastracture/repository/trip"
	location2 "driver_service/internal/infrastracture/services/location"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"log"
	"moul.io/chizap"
	"net/http"
)

type App struct {
	//userHandler    *handlers.UserHandler
	tripService     *trip.Service
	driverService   *driver.Service
	locationService location.LocationService
	repository      domainTrip.Repository

	tripInboundConsumer consumer.Consumer
	tripEventsHandler   consumer.MessageHandler

	tripOutboundProducer trip_outbound2.Producer

	logger *zap.Logger
	srv    *http.Server
	cfg    *config.Config
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(appEnv string) error {
	cfg, err := config.NewConfig(appEnv)
	if err != nil {
		log.Fatal("Could not read config.", err)
		return err
	}
	a.cfg = &cfg

	// Marshal the configuration struct into JSON
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatal("Could not marshal config to JSON.", err)
	}

	// Print the JSON configuration
	fmt.Println("Configuration:")
	fmt.Println(string(jsonCfg))

	a.InitLogger(appEnv)
	a.InitRepositories()
	a.InitProducers()
	a.InitServices()
	a.InitConsumers()

	a.newHttpServer()

	return nil
}

func (a *App) InitLogger(appEnv string) {
	isDebug := appEnv == "development"

	appLogger, err := logger.GetLogger(isDebug)
	if err != nil {
		log.Fatal("Could not initialize logger.", err)
	}
	a.logger = appLogger
}

func (a *App) InitRepositories() {
	a.repository = infraTrip.NewRepository(a.cfg, a.logger)
}

func (a *App) InitServices() {
	a.tripService = trip.NewService(a.repository, a.logger)
	a.locationService = location2.NewService(&http3.Config{Url: a.cfg.Http.LocationServerUrl}, a.logger)
	a.driverService = driver.NewService(a.tripService, a.locationService, a.tripOutboundProducer, a.logger)
}

func (a *App) InitProducers() {
	var async = flag.Bool("a", false, "use async")

	a.tripOutboundProducer = trip_outbound3.NewProducer(async, producer2.Config{
		Host:  a.cfg.Kafka.Host,
		Topic: a.cfg.Kafka.TripInboundTopic,
	}, a.logger)
}

func (a *App) InitConsumers() {
	a.tripEventsHandler = trip_outbound.NewTripInboundMessageHandler(a.tripService, a.logger)
	a.tripInboundConsumer = kafka_consumers.NewConsumer(consumer2.Config{
		Host:           a.cfg.Kafka.Host,
		Topic:          a.cfg.Kafka.TripOutboundTopic,
		Group:          a.cfg.Kafka.TripOutboundGroup,
		SessionTimeout: 6,
		RetryTimeout:   1,
	}, a.tripEventsHandler, a.logger)
}

func (a *App) newHttpServer() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(chizap.New(a.logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	//routers.AddUserRoutes(router, a.userHandler)

	tripServer := http2.NewDriverServer(a.driverService, a.tripService)

	// Создаем экземпляр StrictServer с использованием ранее созданного сервера
	petStoreStrictHandler := openapi.NewStrictHandler(tripServer, nil)
	openapi.HandlerFromMux(petStoreStrictHandler, router)

	a.srv = &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(a.cfg.Http.ServerPort)),
	}
}

func (a *App) Start(ctx context.Context) error {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9000", nil)

	go func() {
		// shutdown := tracing.InitProvider()
		// defer shutdown()

		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Could not listen on port: "+strconv.Itoa(a.cfg.Http.ServerPort), zap.Error(err))
		}
	}()

	go func() {
		a.tripInboundConsumer.Consume(ctx)
		<-ctx.Done()
	}()

	a.logger.Info("Service start at port: " + strconv.Itoa(a.cfg.Http.ServerPort))

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	<-ctx.Done()

	done := make(chan bool)
	a.logger.Info("Server is shutting down...")

	// остановка приложения, gracefully shutdown
	go func() {
		if err := a.srv.Shutdown(context.Background()); err != nil {
			a.logger.Error("Could not gracefully shutdown the userHandler: ", zap.Error(err))
		}

		a.logger.Info("Server stopped ")
		close(done)
	}()

	<-done
	return nil
}
