package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type MessageHandler interface {
	Handle(ctx context.Context, message kafka.Message) error
}
