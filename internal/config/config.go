package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type HttpServer struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type Worker struct {
	ReconclineInterval time.Duration
	JobTimeout         time.Duration
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type GracefulSHutdowner struct {
	GlobalShutdownTimeout      time.Duration
	PerInstanceShutdownTimeout time.Duration
}

type Config struct {
	Env                string
	HttpServer         HttpServer
	Database           Database
	Worker             Worker
	GracefulSHutdowner GracefulSHutdowner
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
		Env: getEnv("ENV"),
		HttpServer: HttpServer{
			Addr: getEnv("HTTP_SERVER_ADDR"),
			ReadTimeout: parseTimeDurationFromEnv(
				"HTTP_SERVER_READ_TIMEOUT",
			),
			WriteTimeout: parseTimeDurationFromEnv(
				"HTTP_SERVER_WRITE_TIMEOUT",
			),
			ShutdownTimeout: parseTimeDurationFromEnv(
				"HTTP_SERVER_SHUTDOWN_TIMEOUT",
			),
		},
		Database: Database{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			Name:     getEnv("DB_NAME"),
		},
		Worker: Worker{
			ReconclineInterval: parseTimeDurationFromEnv(
				"WORKER_RECONCLINE_INTERVAL",
			),
			JobTimeout: parseTimeDurationFromEnv("WORKER_JOB_TIMEOUT"),
		},
		GracefulSHutdowner: GracefulSHutdowner{
			GlobalShutdownTimeout: parseTimeDurationFromEnv(
				"SHUTDOWN_GLOBAL_TIMEOUT",
			),
			PerInstanceShutdownTimeout: parseTimeDurationFromEnv(
				"SHUTDOWN_PER_INSTANCE_TIMEOUT",
			),
		},
	}
}

func parseTimeDurationFromEnv(key string) time.Duration {
	value := getEnv(key)

	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("failed to parse %s, err: %v", key, err)
		return time.Second
	}

	return parsedValue
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Fatalf("environment variable %s not set", key)

	return ""
}
