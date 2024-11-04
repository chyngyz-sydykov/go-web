package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationAddress     string
	ApplicationPort        string
	ApplicationEnvironment string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		ApplicationAddress:     getEnv("APPLICATION_ADDRESS", "/"),
		ApplicationPort:        getEnv("APPLICATION_PORT", "1111"),
		ApplicationEnvironment: getEnv("APPLICATION_ENVIRONMENT", "local"),
	}

	return config, nil
}

func LoadDBConfig() (*DBConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	dbConfig := &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_DATABASE", "database_name"),
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
	}

	return dbConfig, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
