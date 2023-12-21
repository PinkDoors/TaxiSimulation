package main

import (
	"context"
	application "driver_service/app"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctxSys := newSystemContext(ctxWithCancel, 5*time.Second, newLogSystemSignalCallback())

	appEnv, appEnvErr := getEnv("APP_ENV")
	if appEnvErr != nil {
		log.Fatal("No app environment provided")
	}

	app := application.NewApp()

	if err := app.Init(ctxSys, appEnv); err != nil {
		log.Fatal("start app failed")
	}

	if err := app.Start(ctxSys); err != nil {
		log.Fatal("start app failed")
	}

	if err := app.Stop(ctxSys); err != nil {
		log.Fatalf("stop app failed")
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

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}

	return "", errors.New("Unable to get \"" + key + "\" value.")
}

func newLogSystemSignalCallback() Callback {
	return func(signal os.Signal) {
		log.Printf("system signal %d (%s) received, context will be canceled shortly\n", signal, signal.String())
	}
}
