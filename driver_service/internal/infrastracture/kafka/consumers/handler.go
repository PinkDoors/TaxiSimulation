package consumers

import "github.com/segmentio/kafka-go"

type EventHandler interface {
	Handle(message kafka.Message) error
}
