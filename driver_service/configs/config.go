package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Http  HTTP
	Db    Database
	Kafka Kafka
}

type Database struct {
	Uri        string `mapstructure:"MONGO_DB_URI"`
	Timeout    int    `mapstructure:"MONGO_DB_CONNECTION_TIMEOUT_SECONDS"`
	Username   string `mapstructure:"MONGO_DB_USERNAME"`
	Password   string `mapstructure:"MONGO_DB_PASSWORD"`
	AuthSource string `mapstructure:"MONGO_DB_AUTH_SOURCE"`
}

type Kafka struct {
	HOST              string `mapstructure:"KAFKA_HOST"`
	TripInboundTopic  string `mapstructure:"TRIP_INBOUND_TOPIC"`
	TripInboundGroup  string `mapstructure:"TRIP_INBOUND_CONSUMER_GROUP"`
	TripOutboundTopic string `mapstructure:"TRIP_OUTBOUND_TOPIC"`
	TripOutboundGroup string `mapstructure:"TRIP_OUTBOUND_CONSUMER_GROUP"`
}

type HTTP struct {
	PORT int `mapstructure:"SERVER_PORT"`
}

func NewConfig(appEnv string) (Config, error) {
	envFile := fmt.Sprintf(".env.%s", appEnv)
	viper.AddConfigPath(".")
	viper.SetConfigName(envFile)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	dbConfig := getDBConfig()
	kafkaConfig := getKafkaConfig()
	httpConfig := getHttpConfig()

	return Config{
		Http:  httpConfig,
		Db:    dbConfig,
		Kafka: kafkaConfig,
	}, nil
}

func getDBConfig() Database {
	var dbConfig Database

	if err := viper.Unmarshal(&dbConfig); err != nil {
		log.Fatal(err)
	}

	return dbConfig
}

func getKafkaConfig() Kafka {
	var kafkaConfig Kafka

	if err := viper.Unmarshal(&kafkaConfig); err != nil {
		log.Fatal(err)
	}

	return kafkaConfig
}

func getHttpConfig() HTTP {
	var httpConfig HTTP

	if err := viper.Unmarshal(&httpConfig); err != nil {
		log.Fatal(err)
	}

	return httpConfig
}
