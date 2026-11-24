package server

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ticketflow/ticketflow/internal/config"
)

func TestNewServerRoutes(t *testing.T) {
	cfg := &config.Config{
		AppEnv:          "test",
		Port:            "8080",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		IdleTimeout:     10 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	srv := NewServer(cfg, logger)
	if srv == nil {
		t.Fatal("expected server instance, got nil")
	}

	testCases := []struct {
		name       string
		path       string
		method     string
		expectCode int
	}{
		{"Root Healthz", "/healthz", http.MethodGet, http.StatusOK},
		{"Root Livez", "/livez", http.MethodGet, http.StatusOK},
		{"API v1 Health", "/api/v1/health", http.MethodGet, http.StatusOK},
		{"Not Found", "/non-existing-route", http.MethodGet, http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			srv.Engine().ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Errorf("path %s: expected %d, got %d", tc.path, tc.expectCode, w.Code)
			}
		})
	}
}

func TestServerGracefulShutdown(t *testing.T) {
	cfg := &config.Config{
		AppEnv:          "test",
		Port:            "0", // random available test port
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    2 * time.Second,
		IdleTimeout:     5 * time.Second,
		ShutdownTimeout: 2 * time.Second,
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	srv := NewServer(cfg, logger)

	ctx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Run(ctx)
	}()

	// Allow server goroutine to boot
	time.Sleep(100 * time.Millisecond)

	// Trigger shutdown
	cancel()

	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("expected clean shutdown, got error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for graceful shutdown")
	}
}
