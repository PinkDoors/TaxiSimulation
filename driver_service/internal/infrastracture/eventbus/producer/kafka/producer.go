package kafka

import (
	"context"
	"driver_service/internal/infrastracture/eventbus/consumer"
	"flag"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	config  ProducerConfig
	handler consumer.MessageHandler
	logger  *zap.Logger
}

func NewProducer(
	config ProducerConfig,
	handler consumer.MessageHandler,
	logger *zap.Logger,
) *Producer {
	return &Producer{
		config:  config,
		handler: handler,
		logger:  logger,
	}
}

func (kp *Producer) Produce(ctx context.Context, value []byte) {
	go func() {
		var async = flag.Bool("a", false, "use async")

		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers:   []string{"127.0.0.1:29092", "127.0.0.1:39092", "127.0.0.1:49092"},
			Topic:     "demo",
			Async:     *async,
			BatchSize: 10,
		})
		defer writer.Close()

		err := writer.WriteMessages(ctx, kafka.Message{Value: value})
		if err != nil {
			logString := fmt.Sprintf("Consumer for topic: \"%s\" was stopped by cancellation token.", kp.config.Topic)
			kp.logger.Error(logString, zap.Error(err))
		}
	}()

	<-ctx.Done()
}
