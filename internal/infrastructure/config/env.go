package config

import "os"

const (
	defaultEventsBaseURL = "http://localhost:8080"
	defaultEventsAPIKey  = "dev-api-key"
	defaultPostgresDSN   = "postgresql://postgres:root@localhost:5432/users?sslmode=disable"
)

type AppConfig struct {
	EventsBaseURL string
	EventsAPIKey  string
	PostgresDSN   string
}

func Load() AppConfig {
	return AppConfig{
		EventsBaseURL: getEnv("EVENTS_BASE_URL", defaultEventsBaseURL),
		EventsAPIKey:  getEnv("EVENTS_API_KEY", defaultEventsAPIKey),
		PostgresDSN:   getEnv("POSTGRES_DSN", defaultPostgresDSN),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
