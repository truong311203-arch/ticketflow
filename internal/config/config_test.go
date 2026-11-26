package config

import (
	"os"
	"testing"
	"time"
)

func clearConfigEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		"APP_ENV",
		"PORT",
		"SERVER_READ_TIMEOUT",
		"SERVER_WRITE_TIMEOUT",
		"SERVER_IDLE_TIMEOUT",
		"SERVER_SHUTDOWN_TIMEOUT",
		"LOG_LEVEL",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"DB_SCHEMA",
	}
	for _, k := range keys {
		os.Unsetenv(k)
	}
}

func TestLoadDefaults(t *testing.T) {
	clearConfigEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error loading default config, got %v", err)
	}

	if cfg.AppEnv != "development" {
		t.Errorf("expected AppEnv=development, got %s", cfg.AppEnv)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port=8080, got %s", cfg.Port)
	}
	if cfg.ReadTimeout != 10*time.Second {
		t.Errorf("expected ReadTimeout=10s, got %v", cfg.ReadTimeout)
	}
	if cfg.IsProduction() {
		t.Errorf("expected IsProduction()=false for development environment")
	}

	if cfg.DB.Host != "localhost" {
		t.Errorf("expected DB.Host=localhost, got %s", cfg.DB.Host)
	}
	if cfg.DB.Port != "5432" {
		t.Errorf("expected DB.Port=5432, got %s", cfg.DB.Port)
	}
	if cfg.DB.User != "ticketflow_user" {
		t.Errorf("expected DB.User=ticketflow_user, got %s", cfg.DB.User)
	}
	if cfg.DB.Password != "ticketflow_secret" {
		t.Errorf("expected DB.Password=ticketflow_secret, got %s", cfg.DB.Password)
	}
	if cfg.DB.DBName != "ticketflow_db" {
		t.Errorf("expected DB.DBName=ticketflow_db, got %s", cfg.DB.DBName)
	}
	if cfg.DB.SSLMode != "disable" {
		t.Errorf("expected DB.SSLMode=disable, got %s", cfg.DB.SSLMode)
	}
	if cfg.DB.Schema != "core" {
		t.Errorf("expected DB.Schema=core, got %s", cfg.DB.Schema)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	clearConfigEnv(t)

	t.Setenv("APP_ENV", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("SERVER_READ_TIMEOUT", "20s")
	t.Setenv("DB_HOST", "postgres")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "custom_user")
	t.Setenv("DB_PASSWORD", "custom_secret")
	t.Setenv("DB_NAME", "custom_db")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("DB_SCHEMA", "auth")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.AppEnv != "production" {
		t.Errorf("expected AppEnv=production, got %s", cfg.AppEnv)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected Port=9090, got %s", cfg.Port)
	}
	if cfg.ReadTimeout != 20*time.Second {
		t.Errorf("expected ReadTimeout=20s, got %v", cfg.ReadTimeout)
	}
	if !cfg.IsProduction() {
		t.Errorf("expected IsProduction()=true for production environment")
	}

	if cfg.DB.Host != "postgres" {
		t.Errorf("expected DB.Host=postgres, got %s", cfg.DB.Host)
	}
	if cfg.DB.Port != "5433" {
		t.Errorf("expected DB.Port=5433, got %s", cfg.DB.Port)
	}
	if cfg.DB.User != "custom_user" {
		t.Errorf("expected DB.User=custom_user, got %s", cfg.DB.User)
	}
	if cfg.DB.Password != "custom_secret" {
		t.Errorf("expected DB.Password=custom_secret, got %s", cfg.DB.Password)
	}
	if cfg.DB.DBName != "custom_db" {
		t.Errorf("expected DB.DBName=custom_db, got %s", cfg.DB.DBName)
	}
	if cfg.DB.SSLMode != "require" {
		t.Errorf("expected DB.SSLMode=require, got %s", cfg.DB.SSLMode)
	}
	if cfg.DB.Schema != "auth" {
		t.Errorf("expected DB.Schema=auth, got %s", cfg.DB.Schema)
	}
}

func TestValidateInvalidPort(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{"non-numeric", "invalid"},
		{"zero", "0"},
		{"negative", "-1"},
		{"above max", "65536"},
	}

	for _, tc := range testCases {
		t.Run("PORT_"+tc.name, func(t *testing.T) {
			clearConfigEnv(t)
			t.Setenv("PORT", tc.value)
			if _, err := Load(); err == nil {
				t.Errorf("expected error for PORT=%q, got nil", tc.value)
			}
		})

		t.Run("DB_PORT_"+tc.name, func(t *testing.T) {
			clearConfigEnv(t)
			t.Setenv("DB_PORT", tc.value)
			if _, err := Load(); err == nil {
				t.Errorf("expected error for DB_PORT=%q, got nil", tc.value)
			}
		})
	}
}

func TestValidateInvalidEnv(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("APP_ENV", "unknown_env")
	_, err := Load()
	if err == nil {
		t.Errorf("expected error for invalid APP_ENV, got nil")
	}
}

func TestValidateEmptyDBFields(t *testing.T) {
	testCases := []struct {
		name   string
		mutate func(*Config)
	}{
		{"empty Host", func(c *Config) { c.DB.Host = "" }},
		{"empty User", func(c *Config) { c.DB.User = "" }},
		{"empty Password", func(c *Config) { c.DB.Password = "" }},
		{"empty DBName", func(c *Config) { c.DB.DBName = "" }},
		{"empty SSLMode", func(c *Config) { c.DB.SSLMode = "" }},
		{"empty Schema", func(c *Config) { c.DB.Schema = "" }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				AppEnv: "development",
				Port:   "8080",
				DB: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "ticketflow_user",
					Password: "ticketflow_secret",
					DBName:   "ticketflow_db",
					SSLMode:  "disable",
					Schema:   "core",
				},
			}
			tc.mutate(cfg)
			if err := cfg.validate(); err == nil {
				t.Errorf("expected validation error for %s, got nil", tc.name)
			}
		})
	}
}

func TestDBConfigBuildDSN(t *testing.T) {
	t.Run("standard credentials", func(t *testing.T) {
		dbCfg := DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "ticketflow_user",
			Password: "ticketflow_secret",
			DBName:   "ticketflow_db",
			SSLMode:  "disable",
			Schema:   "core",
		}

		expected := "postgres://ticketflow_user:ticketflow_secret@localhost:5432/ticketflow_db?search_path=core&sslmode=disable"
		if got := dbCfg.BuildDSN(); got != expected {
			t.Errorf("expected DSN %q, got %q", expected, got)
		}
	})

	t.Run("password with special characters is URL encoded", func(t *testing.T) {
		dbCfg := DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "ticketflow_user",
			Password: "p@ss#word/1",
			DBName:   "ticketflow_db",
			SSLMode:  "disable",
			Schema:   "core",
		}

		expected := "postgres://ticketflow_user:p%40ss%23word%2F1@localhost:5432/ticketflow_db?search_path=core&sslmode=disable"
		if got := dbCfg.BuildDSN(); got != expected {
			t.Errorf("expected DSN %q, got %q", expected, got)
		}
	})
}
