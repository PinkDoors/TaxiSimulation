package kafka

import (
	"context"
	consumer2 "driver_service/configs/kafka/consumer"
	"driver_service/internal/infrastracture/eventbus/consumer"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type Consumer struct {
	config  consumer2.Config
	handler consumer.MessageHandler
	logger  *zap.Logger
}

func NewConsumer(
	config consumer2.Config,
	handler consumer.MessageHandler,
	logger *zap.Logger,
) *Consumer {
	return &Consumer{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

func (c *Consumer) Consume(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{c.config.Host},
		Topic:          c.config.Topic,
		GroupID:        c.config.Group,
		SessionTimeout: time.Second * time.Duration(c.config.SessionTimeout),
	})
	defer reader.Close()

	logString := fmt.Sprintf("Started consumer for topic: \"%s\".", c.config.Topic)
	c.logger.Info(logString)

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logString := fmt.Sprintf("Consumer for topic: \"%s\" was stopped by cancellation token.", c.config.Topic)
				c.logger.Error(logString, zap.Error(err))
				return
			}
			logString := fmt.Sprintf("Failed to read messages from topic: \"%s\".", c.config.Topic)
			c.logger.Error(logString, zap.Error(err))
			continue
		}

		logString := fmt.Sprintf("Consumed message from topic \"%s\" - \"%s\"", c.config.Topic, string(msg.Value))
		c.logger.Info(logString)
		err = c.handler.Handle(ctx, msg)

		if err != nil {
			logString := fmt.Sprintf("Error trying to handle message from topic: \"%s\".", c.config.Topic)
			c.logger.Error(logString, zap.Error(err))
		} else {
			err := reader.CommitMessages(ctx, msg)
			logString := fmt.Sprintf("Commmited message from topic \"%s\" - \"%s\"", c.config.Topic, string(msg.Value))
			c.logger.Info(logString)
			if err != nil {
				logString := fmt.Sprintf("Error trying to commit message from topic: \"%s\".", c.config.Topic)
				c.logger.Error(logString, zap.Error(err))
			}
		}

		time.Sleep(300 * time.Duration(c.config.RetryTimeout))
	}
}
