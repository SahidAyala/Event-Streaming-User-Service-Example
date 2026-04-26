package main

import "os"

const (
	defaultEventsBaseURL = "http://localhost:8080"
	defaultEventsAPIKey  = "dev-api-key"
	defaultPostgresDSN   = "postgresql://postgres:root@localhost:5432/users"
)

type config struct {
	eventsBaseURL string
	eventsAPIKey  string
	postgresDSN   string
}

func loadConfig() config {
	return config{
		eventsBaseURL: getEnv("EVENTS_BASE_URL", defaultEventsBaseURL),
		eventsAPIKey:  getEnv("EVENTS_API_KEY", defaultEventsAPIKey),
		postgresDSN:   getEnv("POSTGRES_DSN", defaultPostgresDSN),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
