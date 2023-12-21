package kafka

type ConsumerConfig struct {
	Host           string
	Topic          string
	Group          string
	SessionTimeout int
	RetryTimeout   int
}
