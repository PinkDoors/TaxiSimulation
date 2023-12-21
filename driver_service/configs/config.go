package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv string
	PORT   int
	Db     Database
	Kafka  Kafka
}

type Database struct {
	Uri        string
	Timeout    int
	Username   string
	Password   string
	Authsource string
}

type Kafka struct {
	TripInboundTopic    string
	TripInboundGroupId  string
	TripOutboundTopic   string
	TripOutboundGroupId string
}

type HTTP struct {
	PORT string
}

func NewConfig() (*Config, error) {
	appEnv, appEnvErr := getEnv("APP_ENV")
	if appEnvErr != nil {
		return nil, appEnvErr
	}

	envFile := fmt.Sprintf(".env.%s", appEnv)
	if err := godotenv.Load(envFile); err != nil {
		return nil, err
	}

	port, portErr := getEnvAsInt("PORT")
	if portErr != nil {
		return nil, portErr
	}

	db, dbErr := getDBConfig()
	if dbErr != nil {
		return nil, dbErr
	}

	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName("app")

	envFile := fmt.Sprintf(".env.%s", appEnv)
	// Tell viper the type of your file
	viper.SetConfigType(envFile)

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return &Config{
		AppEnv: appEnv,
		PORT:   port,
		Db:     db,
	}, nil
}

func getDBConfig() (Database, error) {
	uri, uriErr := getEnv("MONGO_DB_URI")
	if uriErr != nil {
		return Database{}, uriErr
	}

	timeout, timeoutErr := getEnvAsInt("MONGO_DB_CONNECTION_TIMEOUT_SECONDS")
	if timeoutErr != nil {
		return Database{}, timeoutErr
	}

	username, usernameErr := getEnv("MONGO_DB_USERNAME")
	if usernameErr != nil {
		return Database{}, usernameErr
	}

	password, passwordErr := getEnv("MONGO_DB_PASSWORD")
	if passwordErr != nil {
		return Database{}, passwordErr
	}

	authSource, authSourceErr := getEnv("MONGO_DB_AUTH_SOURCE")
	if passwordErr != nil {
		return Database{}, authSourceErr
	}

	return Database{
		Uri:        uri,
		Timeout:    timeout,
		Username:   username,
		Password:   password,
		Authsource: authSource,
	}, nil
}

func getKafkaConfig() (Kafka, error) {
	tripInboundTopic, tripInboundTopicErr := getEnv("TRIP_INBOUND_TOPIC")
	if tripInboundTopicErr != nil {
		return Kafka{}, tripInboundTopicErr
	}

	tripInboundConsumerGroup, tripInboundConsumerGroupErr := getEnv("TRIP_INBOUND_CONSUMER_GROUP")
	if tripInboundConsumerGroupErr != nil {
		return Kafka{}, tripInboundConsumerGroupErr
	}

	tripOutboundTopic, tripOutboundTopicErr := getEnv("TRIP_INBOUND_TOPIC")
	if tripOutboundTopicErr != nil {
		return Kafka{}, tripOutboundTopicErr
	}

	tripOutboundConsumerGroup, tripOutboundConsumerGroupErr := getEnv("TRIP_INBOUND_CONSUMER_GROUP")
	if tripOutboundConsumerGroupErr != nil {
		return Kafka{}, tripOutboundConsumerGroupErr
	}

	return Kafka{
		TripInboundTopic:    tripInboundTopic,
		TripInboundGroupId:  tripInboundConsumerGroup,
		TripOutboundTopic:   tripOutboundTopic,
		TripOutboundGroupId: tripOutboundConsumerGroup,
	}, nil
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}

	return "", errors.New("Unable to get \"" + key + "\" value.")
}

func getEnvAsInt(key string) (int, error) {
	valueStr, parseErr := getEnv(key)
	if parseErr != nil {
		return 0, parseErr
	}

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value, nil
	}

	return 0, errors.New("Unable to parse \"" + key + "\" value to Int.")
}
