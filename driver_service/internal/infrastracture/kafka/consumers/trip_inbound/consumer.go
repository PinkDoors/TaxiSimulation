package trip_inbound

import (
	"context"
	"driver_service/internal/infrastracture/kafka/consumers"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TripBoundConsumer struct {
	config  *consumers.ConsumerConfig
	handler consumers.EventHandler
	logger  *zap.Logger
}

func NewRepository(
	config *consumers.ConsumerConfig,
	handler consumers.EventHandler,
	logger *zap.Logger,
) *TripBoundConsumer {
	return &TripBoundConsumer{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

func (tbc *TripBoundConsumer) Consume(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	go func() {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{tbc.config.Host},
			Topic:          tbc.config.Topic,
			GroupID:        tbc.config.Group,
			SessionTimeout: time.Second * time.Duration(tbc.config.SessionTimeout),
		})
		defer reader.Close()

		for {
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				logString := fmt.Sprintf("Failed to read messages from topic: %s", tbc.config.Topic)
				tbc.logger.Error(logString, zap.Error(err))
			}

			err = tbc.handler.Handle(msg)

			if err != nil {
				logString := fmt.Sprintf("Error trying to handle %s message.", tbc.config.Topic)
				tbc.logger.Error(logString, zap.Error(err))
			} else {
				err := reader.CommitMessages(ctx, msg)
				if err != nil {
					logString := fmt.Sprintf("Error trying to commit %s message.", tbc.config.Topic)
					tbc.logger.Error(logString, zap.Error(err))
				}
			}

			time.Sleep(300 * time.Duration(tbc.config.RetryTimeout))
		}
	}()

	<-ctx.Done()
}
