package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ticketflow/ticketflow/internal/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockPinger struct {
	err error
}

func (m *mockPinger) Ping(_ context.Context) error {
	return m.err
}

func TestHealthHandler_Healthz(t *testing.T) {
	cfg := &config.Config{
		AppEnv: "test",
		Port:   "8080",
	}
	h := NewHealthHandler(cfg, &mockPinger{})

	router := gin.New()
	router.GET("/healthz", h.Healthz)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", res.Status)
	}
	if res.Application != "ticketflow" {
		t.Errorf("expected application 'ticketflow', got %q", res.Application)
	}
	if res.Environment != "test" {
		t.Errorf("expected environment 'test', got %q", res.Environment)
	}
	if res.Database != "up" {
		t.Errorf("expected database 'up', got %q", res.Database)
	}
}

func TestHealthHandler_Healthz_DBDown(t *testing.T) {
	cfg := &config.Config{
		AppEnv: "test",
		Port:   "8080",
	}

	testCases := []struct {
		name string
		db   DBPinger
	}{
		{"ping error", &mockPinger{err: errors.New("connection refused")}},
		{"nil db", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHealthHandler(cfg, tc.db)

			router := gin.New()
			router.GET("/healthz", h.Healthz)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusServiceUnavailable {
				t.Fatalf("expected status 503, got %d", w.Code)
			}

			var res HealthResponse
			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if res.Status != "degraded" {
				t.Errorf("expected status 'degraded', got %q", res.Status)
			}
			if res.Database != "down" {
				t.Errorf("expected database 'down', got %q", res.Database)
			}
		})
	}
}

func TestHealthHandler_Livez(t *testing.T) {
	cfg := &config.Config{AppEnv: "test"}
	// Even when DB is down, liveness probe must still return 200 OK.
	h := NewHealthHandler(cfg, &mockPinger{err: errors.New("db down")})

	router := gin.New()
	router.GET("/livez", h.Livez)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/livez", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res["status"] != "alive" {
		t.Errorf("expected status 'alive', got %q", res["status"])
	}
}
