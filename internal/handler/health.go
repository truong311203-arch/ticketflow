package handler

import (
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
}

// HealthHandler handles health check requests.
type HealthHandler struct {
	cfg *config.Config
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

// Healthz responds with a full system status probe.
func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:      "ok",
		Application: "ticketflow",
		Environment: h.cfg.AppEnv,
		Timestamp:   time.Now().UTC(),
	})
}

// Livez responds with a lightweight liveness probe.
func (h *HealthHandler) Livez(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
