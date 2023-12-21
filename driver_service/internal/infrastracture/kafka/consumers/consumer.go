package consumers

import "context"

type Consumer interface {
	Consume(context.Context) error
}
