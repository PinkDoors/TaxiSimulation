package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	AppEnv string
	PORT   int
	DB     Database
}

type Database struct {
	URI        string
	TIMEOUT    int
	USERNAME   string
	PASSWORD   string
	AUTHSOURCE string
}

type HTTP struct {
	PORT string
}

func NewConfig() (*Config, error) {
	appEnv, appEnvErr := getEnv("APP_ENV")
	if appEnvErr != nil {
		return nil, appEnvErr
	}

	port, portErr := getEnvAsInt("PORT")
	if portErr != nil {
		return nil, portErr
	}

	db, dbErr := getDBConfig()
	if dbErr != nil {
		return nil, dbErr
	}

	return &Config{
		AppEnv: appEnv,
		PORT:   port,
		DB:     db,
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
		URI:        uri,
		TIMEOUT:    timeout,
		USERNAME:   username,
		PASSWORD:   password,
		AUTHSOURCE: authSource,
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
