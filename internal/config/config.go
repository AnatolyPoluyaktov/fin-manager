package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Env            string
	AuthToken      string
	Storage        Storage
	Server         HTTPServer
	MigrationsPath string
}

type Storage struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

type HTTPServer struct {
	Address string
	Timeout time.Duration
}

func MustLoadConfig() *Config {
	config := &Config{
		Env:       getEnv("ENV", "local"),
		AuthToken: mustGetEnv("AUTH_TOKEN"),
		Storage: Storage{
			DbUser:     mustGetEnv("DB_USER"),
			DbPassword: mustGetEnv("DB_PASSWORD"),
			DbHost:     mustGetEnv("DB_HOST"),
			DbPort:     mustGetEnv("DB_PORT"),
			DbName:     mustGetEnv("DB_NAME"),
		},
		Server: HTTPServer{
			Address: getEnv("SERVER_ADDRESS", "localhost:8080"),
			Timeout: getEnvAsDuration("SERVER_TIMEOUT", 4*time.Second),
		},
		MigrationsPath: getEnv("MIGRATIONS_PATH", "migrations"),
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		d, err := time.ParseDuration(value)
		if err != nil {
			log.Fatalf("Invalid duration format for %s: %v", key, err)
		}
		return d
	}
	return defaultValue
}
