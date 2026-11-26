package config

import (
	"fmt"
	"net"
	"net/url"
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
	DB              DBConfig
}

// DBConfig holds PostgreSQL connection parameters.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Schema   string
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
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "ticketflow_user"),
			Password: getEnv("DB_PASSWORD", "ticketflow_secret"),
			DBName:   getEnv("DB_NAME", "ticketflow_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Schema:   getEnv("DB_SCHEMA", "core"),
		},
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

// BuildDSN builds a URL-encoded PostgreSQL connection string for pgxpool.
func (d DBConfig) BuildDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   net.JoinHostPort(d.Host, d.Port),
		Path:   d.DBName,
	}

	q := u.Query()
	q.Set("sslmode", d.SSLMode)
	q.Set("search_path", d.Schema)
	u.RawQuery = q.Encode()

	return u.String()
}

func validatePort(name, val string) error {
	if val == "" {
		return fmt.Errorf("%s must not be empty", name)
	}

	portNum, err := strconv.Atoi(val)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("%s must be a valid port number between 1 and 65535, got: %s", name, val)
	}

	return nil
}

// validate ensures core configuration values are valid.
func (c *Config) validate() error {
	if err := validatePort("PORT", c.Port); err != nil {
		return err
	}
	if err := validatePort("DB_PORT", c.DB.Port); err != nil {
		return err
	}

	requiredDBFields := []struct {
		name  string
		value string
	}{
		{"DB_HOST", c.DB.Host},
		{"DB_USER", c.DB.User},
		{"DB_PASSWORD", c.DB.Password},
		{"DB_NAME", c.DB.DBName},
		{"DB_SSLMODE", c.DB.SSLMode},
		{"DB_SCHEMA", c.DB.Schema},
	}

	for _, field := range requiredDBFields {
		if field.value == "" {
			return fmt.Errorf("%s must not be empty", field.name)
		}
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
