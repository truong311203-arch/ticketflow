package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	// Ensure clean environment for test
	os.Unsetenv("APP_ENV")
	os.Unsetenv("PORT")
	os.Unsetenv("SERVER_READ_TIMEOUT")
	os.Unsetenv("SERVER_WRITE_TIMEOUT")
	os.Unsetenv("SERVER_IDLE_TIMEOUT")
	os.Unsetenv("SERVER_SHUTDOWN_TIMEOUT")
	os.Unsetenv("LOG_LEVEL")

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
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("SERVER_READ_TIMEOUT", "20s")

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
}

func TestValidateInvalidPort(t *testing.T) {
	t.Setenv("PORT", "invalid")
	_, err := Load()
	if err == nil {
		t.Errorf("expected error for invalid port, got nil")
	}
}

func TestValidateInvalidEnv(t *testing.T) {
	t.Setenv("APP_ENV", "unknown_env")
	_, err := Load()
	if err == nil {
		t.Errorf("expected error for invalid APP_ENV, got nil")
	}
}
