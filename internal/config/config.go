package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	// GitHub personal access token
	GitHubToken string

	// Server configuration
	ServerPort string
	ServerHost string

	// OAuth configuration (for future use)
	OAuthClientID     string
	OAuthClientSecret string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}

	cfg := &Config{
		GitHubToken:       token,
		ServerPort:        getEnvOrDefault("SERVER_PORT", "8080"),
		ServerHost:        getEnvOrDefault("SERVER_HOST", "localhost"),
		OAuthClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		OAuthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
