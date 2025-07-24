package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port int
	DB   DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func New() *Config {
	return &Config{
		Port: getEnvAsInt("PORT", 8123),
		DB: DBConfig{
			Host:     getEnv("DATABASE_HOST", "127.0.0.1"),
			Port:     3306,
			Database: "users",
			User:     "users_service",
			Password: "123",
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
