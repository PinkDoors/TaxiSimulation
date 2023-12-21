package kafka

import (
	"context"
	"driver_service/internal/infrastracture/eventbus/consumer"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type Consumer struct {
	config  ConsumerConfig
	handler consumer.MessageHandler
	logger  *zap.Logger
}

func NewKafkaConsumer(
	config ConsumerConfig,
	handler consumer.MessageHandler,
	logger *zap.Logger,
) *Consumer {
	return &Consumer{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

func (tbc *Consumer) Consume(ctx context.Context) {
	go func() {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{tbc.config.Host},
			Topic:          tbc.config.Topic,
			GroupID:        tbc.config.Group,
			SessionTimeout: time.Second * time.Duration(tbc.config.SessionTimeout),
		})
		defer reader.Close()

		logString := fmt.Sprintf("Started consumer for topic: \"%s\".", tbc.config.Topic)
		tbc.logger.Info(logString)

		for {
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					logString := fmt.Sprintf("Consumer for topic: \"%s\" was stopped by cancellation token.", tbc.config.Topic)
					tbc.logger.Error(logString, zap.Error(err))
					return
				}
				logString := fmt.Sprintf("Failed to read messages from topic: \"%s\".", tbc.config.Topic)
				tbc.logger.Error(logString, zap.Error(err))
				continue
			}

			tbc.logger.Info("Message:" + string(msg.Value))
			// err = tbc.handler.Handle(ctx, msg)

			if err != nil {
				logString := fmt.Sprintf("Error trying to handle message from topic: \"%s\".", tbc.config.Topic)
				tbc.logger.Error(logString, zap.Error(err))
			} else {
				err := reader.CommitMessages(ctx, msg)
				if err != nil {
					logString := fmt.Sprintf("Error trying to commit message from topic: \"%s\".", tbc.config.Topic)
					tbc.logger.Error(logString, zap.Error(err))
				}
			}

			time.Sleep(300 * time.Duration(tbc.config.RetryTimeout))
		}
	}()

	<-ctx.Done()
}
