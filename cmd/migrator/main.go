package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(key + " is not set")
}

func main() {
	configPath := getEnv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config file does not exist: " + configPath)
	}

	if err := godotenv.Load(configPath); err != nil {
		panic("Failed to load config: " + err.Error())
	}

	databaseConfig := &databaseConfig{
		Host:     getEnv("DB_HOST"),
		Port:     getEnv("DB_PORT"),
		Username: getEnv("DB_USER"),
		Password: getEnv("DB_PASSWORD"),
		Database: getEnv("DB_NAME"),
	}

	var migrationPath string
	flag.StringVar(&migrationPath, "path", "migrations", "path to the migrations directory")
	flag.Parse()

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		databaseConfig.Username, databaseConfig.Password,
		databaseConfig.Host, databaseConfig.Port, databaseConfig.Database)


	m, err := migrate.New("file://"+migrationPath, dbURL)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		}

		panic(err)
	}

	fmt.Println("Migrations applied successfully")
}

type databaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}