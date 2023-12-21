package kafka

type ProducerConfig struct {
	Host           string
	Topic          string
	Group          string
	SessionTimeout int
	RetryTimeout   int
}
