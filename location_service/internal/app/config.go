package app

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	ShutdownTimeout int
	Port            int
	BasePath        string
}

type DatabaseConfig struct {
	DSN           string
	MigrationsDir string
}

type Config struct {
	App AppConfig
	Db  DatabaseConfig
}

func NewConfig(env string) (*Config, error) {
	filename := ".env." + env
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConf, dbErr := getDbConfig()
	if dbErr != nil {
		return nil, dbErr
	}

	appConfig, appConfigErr := getAppConfig()
	if appConfigErr != nil {
		return nil, appConfigErr
	}

	return &Config{
		App: appConfig,
		Db:  dbConf,
	}, nil
}

func getDbConfig() (DatabaseConfig, error) {
	dbDSN, dbDsnError := getEnv("DSN")
	if dbDsnError != nil {
		return DatabaseConfig{}, dbDsnError
	}

	dbMigrationsDir, dbMigrationsDirError := getEnv("MIGRATION_DIR")
	if dbMigrationsDirError != nil {
		return DatabaseConfig{}, dbMigrationsDirError
	}

	return DatabaseConfig{
		DSN:           dbDSN,
		MigrationsDir: dbMigrationsDir,
	}, nil
}

func getAppConfig() (AppConfig, error) {
	shutdownTimeoutStr, shutdownTimeoutReadError := getEnv("SHUTDOWN_TIMEOUT")
	if shutdownTimeoutReadError != nil {
		return AppConfig{}, shutdownTimeoutReadError
	}

	shutdownTimeout, parseErr := strconv.Atoi(shutdownTimeoutStr)
	if parseErr != nil {
		return AppConfig{}, parseErr
	}

	portStr, portReadError := getEnv("HOST")
	if portReadError != nil {
		return AppConfig{}, portReadError
	}

	port, parseErr := strconv.Atoi(portStr)
	if parseErr != nil {
		return AppConfig{}, parseErr
	}

	basePath, basePathError := getEnv("BASE_PATH")
	if basePathError != nil {
		return AppConfig{}, basePathError
	}

	return AppConfig{
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
		BasePath:        basePath,
	}, nil
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}

	return "", errors.New("Unable to get \"" + key + "\" value.")
}
