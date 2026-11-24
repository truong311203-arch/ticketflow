package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestStructuredLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	r := gin.New()
	r.Use(StructuredLogger(logger))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	logOutput := buf.String()
	if !bytes.Contains([]byte(logOutput), []byte("HTTP request")) {
		t.Errorf("expected log output to contain 'HTTP request', got: %s", logOutput)
	}
	if !bytes.Contains([]byte(logOutput), []byte("/ping")) {
		t.Errorf("expected log output to contain '/ping', got: %s", logOutput)
	}
}

func TestRecovery(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	r := gin.New()
	r.Use(Recovery(logger))
	r.GET("/panic", func(c *gin.Context) {
		panic("unexpected runtime error")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 on panic, got %d", w.Code)
	}

	logOutput := buf.String()
	if !bytes.Contains([]byte(logOutput), []byte("Panic recovered")) {
		t.Errorf("expected log output to contain 'Panic recovered', got: %s", logOutput)
	}
}
