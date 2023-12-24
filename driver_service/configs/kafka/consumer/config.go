package consumer

type Config struct {
	Host           string
	Topic          string
	Group          string
	SessionTimeout int
	RetryTimeout   int
}
