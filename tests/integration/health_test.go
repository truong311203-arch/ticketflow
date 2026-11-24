package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ticketflow/ticketflow/internal/config"
	"github.com/ticketflow/ticketflow/internal/handler"
	"github.com/ticketflow/ticketflow/internal/server"
)

func TestHealthEndpointsE2E(t *testing.T) {
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
	srv := server.NewServer(cfg, logger)

	// Spin up a live local HTTP test server
	ts := httptest.NewServer(srv.Engine())
	defer ts.Close()

	client := ts.Client()

	t.Run("GET /healthz returns 200 and valid JSON", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/healthz")
		if err != nil {
			t.Fatalf("failed to execute GET /healthz: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		var payload handler.HealthResponse
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if payload.Status != "ok" {
			t.Errorf("expected status 'ok', got %q", payload.Status)
		}
		if payload.Application != "ticketflow" {
			t.Errorf("expected app 'ticketflow', got %q", payload.Application)
		}
		if payload.Environment != "test" {
			t.Errorf("expected env 'test', got %q", payload.Environment)
		}
	})

	t.Run("GET /livez returns 200 alive", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/livez")
		if err != nil {
			t.Fatalf("failed to execute GET /livez: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		var payload map[string]string
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if payload["status"] != "alive" {
			t.Errorf("expected status 'alive', got %q", payload["status"])
		}
	})

	t.Run("GET /api/v1/health returns 200", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/api/v1/health")
		if err != nil {
			t.Fatalf("failed to execute GET /api/v1/health: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})
}
