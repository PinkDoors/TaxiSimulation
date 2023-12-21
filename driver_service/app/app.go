package app

import (
	"context"
	config "driver_service/configs"
	http2 "driver_service/internal/adapters/http"
	openapi "driver_service/internal/adapters/http/generate"
	"driver_service/internal/application/trip"
	domainTrip "driver_service/internal/domain/repository/trip"
	"driver_service/internal/infrastracture/eventbus/consumer"
	"driver_service/internal/infrastracture/eventbus/consumer/kafka"
	"driver_service/internal/infrastracture/eventbus/consumer/kafka/trip_inbound"
	"driver_service/internal/infrastracture/logger"
	infraTrip "driver_service/internal/infrastracture/repository/trip"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"log"
	"moul.io/chizap"
	"net/http"
)

type App struct {
	//userHandler    *handlers.UserHandler
	service    *trip.Service
	repository domainTrip.Repository

	tripEventsConsumer consumer.Consumer
	tripEventsHandler  consumer.MessageHandler
	logger             *zap.Logger
	srv                *http.Server
	cfg                *config.Config
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(ctx context.Context, appEnv string) error {
	cfg, err := config.NewConfig(appEnv)
	if err != nil {
		log.Fatal("Could not read config.", err)
		return err
	}
	a.cfg = &cfg

	a.InitLogger(appEnv)
	a.InitRepositories()
	a.InitServices()
	a.InitKafka()

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
	a.service = trip.NewService(a.repository, a.logger)
}

func (a *App) InitKafka() {
	a.tripEventsHandler = trip_inbound.NewTripInboundMessageHandler(a.service, a.logger)
	a.tripEventsConsumer = kafka.NewKafkaConsumer(kafka.ConsumerConfig{
		Host:           a.cfg.Kafka.TripInboundTopic,
		Topic:          a.cfg.Kafka.TripInboundTopic,
		Group:          a.cfg.Kafka.TripInboundGroup,
		SessionTimeout: 1,
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

	tripServer := http2.NewTripServer(*a.service)

	// Создаем экземпляр StrictServer с использованием ранее созданного сервера
	petStoreStrictHandler := openapi.NewStrictHandler(tripServer, nil)
	openapi.HandlerFromMux(petStoreStrictHandler, router)

	a.srv = &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(a.cfg.Http.PORT)),
	}
}

func (a *App) Start(ctx context.Context) error {
	//http.Handle("/metrics", promhttp.Handler())
	//go http.ListenAndServe(":9000", nil) //nolint:errcheck

	go func() {
		// shutdown := tracing.InitProvider()
		// defer shutdown()

		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Could not listen on port: "+strconv.Itoa(a.cfg.Http.PORT), zap.Error(err))
		}
	}()

	go func() {
		a.tripEventsConsumer.Consume(ctx)
	}()

	a.logger.Info("Service start at port: " + strconv.Itoa(a.cfg.Http.PORT))

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
