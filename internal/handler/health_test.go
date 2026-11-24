package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ticketflow/ticketflow/internal/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthHandler_Healthz(t *testing.T) {
	cfg := &config.Config{
		AppEnv: "test",
		Port:   "8080",
	}
	h := NewHealthHandler(cfg)

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
}

func TestHealthHandler_Livez(t *testing.T) {
	cfg := &config.Config{AppEnv: "test"}
	h := NewHealthHandler(cfg)

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
