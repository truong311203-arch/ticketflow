package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ticketflow/ticketflow/internal/config"
)

// HealthResponse represents the payload returned by the health check endpoint.
type HealthResponse struct {
	Status      string    `json:"status"`
	Application string    `json:"application"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
	Database    string    `json:"database"`
}

type DBPinger interface {
	Ping(ctx context.Context) error
}

// HealthHandler handles health check requests.
type HealthHandler struct {
	cfg *config.Config
	db  DBPinger
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(cfg *config.Config, db DBPinger) *HealthHandler {
	return &HealthHandler{cfg: cfg, db: db}
}

// Healthz responds with a full system status probe.
func (h *HealthHandler) Healthz(c *gin.Context) {
	pingCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if h.db == nil || h.db.Ping(pingCtx) != nil {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status:      "degraded",
			Application: "ticketflow",
			Environment: h.cfg.AppEnv,
			Timestamp:   time.Now().UTC(),
			Database:    "down",
		})
		return
	}
	c.JSON(http.StatusOK, HealthResponse{
		Status:      "ok",
		Application: "ticketflow",
		Environment: h.cfg.AppEnv,
		Timestamp:   time.Now().UTC(),
		Database:    "up",
	})
}

// Livez responds with a lightweight liveness probe.
func (h *HealthHandler) Livez(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
