package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBConnectionString string
	ServerPort         string
	SlackClientID      string
	SlackClientSecret  string
}

var Config AppConfig

func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Config = AppConfig{
		DBConnectionString: getEnv("DB_CONNECTION_STRING", ""),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		SlackClientID:      getEnv("SLACK_CLIENT_ID", ""),
		SlackClientSecret:  getEnv("SLACK_CLIENT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}