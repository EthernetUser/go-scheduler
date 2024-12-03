package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)


type HttpServer struct {
	Addr string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Host string
	Port string
	User string
	Password string
	Name string
}

type Config struct {
	Env string
	HttpServer HttpServer
	Database Database
}



func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("failed to load config file %s, err: %v", configPath, err)
	}

	return &Config{
		Env: getEnv("ENV", "local"),
		HttpServer: HttpServer{
			Addr: getEnv("HTTP_SERVER_ADDR", ":8080"),
			ReadTimeout: parseTimeDurationFromEnv("HTTP_SERVER_READ_TIMEOUT", "10s"),
			WriteTimeout: parseTimeDurationFromEnv("HTTP_SERVER_WRITE_TIMEOUT", "10s"),
		},
		Database: Database{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name: getEnv("DB_NAME", "postgres"),
		},
	}
}


func parseTimeDurationFromEnv(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)

	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("failed to parse %s, err: %v", key, err)
	}

	return parsedValue
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultValue == "" {
		log.Fatalf("environment variable %s not set", key)
	}

	return defaultValue
}