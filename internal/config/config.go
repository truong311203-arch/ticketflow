package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration parameters.
type Config struct {
	AppEnv          string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	LogLevel        string
}

// Load loads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:     getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 15*time.Second),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// IsProduction returns true if the current environment is production.
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// validate ensures core configuration values are valid.
func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT must not be empty")
	}

	portNum, err := strconv.Atoi(c.Port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("PORT must be a valid port number between 1 and 65535, got: %s", c.Port)
	}

	switch c.AppEnv {
	case "development", "staging", "production", "test":
		// valid
	default:
		return fmt.Errorf("APP_ENV must be one of [development, staging, production, test], got: %s", c.AppEnv)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}
