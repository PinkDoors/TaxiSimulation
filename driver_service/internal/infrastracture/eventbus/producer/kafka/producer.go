package kafka

import (
	"context"
	"driver_service/configs/kafka/producer"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	async  *bool
	config producer.Config
	logger *zap.Logger
}

func NewProducer(
	async *bool,
	config producer.Config,
	logger *zap.Logger,
) *Producer {
	return &Producer{
		async:  async,
		config: config,
		logger: logger,
	}
}

func (p *Producer) Produce(ctx context.Context, value []byte) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{p.config.Host},
		Topic:     p.config.Topic,
		Async:     *p.async,
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
