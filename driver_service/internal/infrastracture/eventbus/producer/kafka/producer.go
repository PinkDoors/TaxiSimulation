package kafka

import (
	"context"
	"flag"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	config ProducerConfig
	logger *zap.Logger
}

func NewProducer(
	config ProducerConfig,
	logger *zap.Logger,
) *Producer {
	return &Producer{
		config: config,
		logger: logger,
	}
}

func (p *Producer) Produce(ctx context.Context, value []byte) error {
	var async = flag.Bool("a", false, "use async")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{p.config.Host},
		Topic:     p.config.Topic,
		Async:     *async,
		BatchSize: 10,
	})
	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{Value: value})
	if err != nil {
		logString := fmt.Sprintf("Consumer for topic: \"%s\" was stopped by cancellation token.", p.config.Topic)
		p.logger.Error(logString, zap.Error(err))
		return err
	}

	return nil
}
