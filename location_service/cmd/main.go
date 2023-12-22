package main

import (
	"context"
	"location_service/internal/app"
	"location_service/internal/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := logger.GetLogger(false)
	if err != nil {
		log.Fatal(err)
	}

	env, exists := os.LookupEnv("APP_ENV")
	if exists != true {
		log.Fatal("No APP_ENV variable")
	}

	config, err := app.NewConfig(env)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctxSys := newSystemContext(ctxWithCancel, time.Duration(config.App.ShutdownTimeout), newLogSystemSignalCallback())

	a, err := app.New(config, logger)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := a.Serve(); err != nil {
		log.Fatal(err.Error())
	}

	if err := a.Shutdown(ctxSys); err != nil {
		log.Fatalf(err.Error())
	}
}

type Callback func(signal os.Signal)

// NewSystemContext returns new Context, which will be cancelled on receiving SIGTERM and SIGINT signals after supplied delay.
// Additionally multiple Callback functions can be passed, they will be called immediately after receiving signals, before delay.
func newSystemContext(ctx context.Context, delay time.Duration, callbacks ...Callback) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		sig := <-sigint
		for _, cb := range callbacks {
			go cb(sig)
		}

		time.Sleep(delay)

		cancel()
	}()

	return ctx
}

func newLogSystemSignalCallback() Callback {
	return func(signal os.Signal) {
		log.Printf("system signal %d (%s) received, context will be canceled shortly\n", signal, signal.String())
	}
}
