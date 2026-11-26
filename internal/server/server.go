package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ticketflow/ticketflow/internal/config"
	"github.com/ticketflow/ticketflow/internal/handler"
	"github.com/ticketflow/ticketflow/internal/middleware"
)

// Server wraps the HTTP server and engine.
type Server struct {
	httpServer *http.Server
	engine     *gin.Engine
	cfg        *config.Config
	logger     *slog.Logger
}

// NewServer initializes the HTTP server, engine, middlewares, and routes.
func NewServer(cfg *config.Config, logger *slog.Logger, db handler.DBPinger) *Server {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.AppEnv == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(middleware.StructuredLogger(logger))
	engine.Use(middleware.Recovery(logger))

	// Register root health and readiness probes
	healthHdl := handler.NewHealthHandler(cfg, db)
	engine.GET("/healthz", healthHdl.Healthz)
	engine.GET("/livez", healthHdl.Livez)

	// API versioning group (extensible for future business modules)
	v1 := engine.Group("/api/v1")
	{
		v1.GET("/health", healthHdl.Healthz)
	}

	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpSrv,
		engine:     engine,
		cfg:        cfg,
		logger:     logger,
	}
}

// Engine returns the underlying gin.Engine (useful for testing).
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Run starts the server and listens for context cancellation for graceful shutdown.
func (s *Server) Run(ctx context.Context) error {
	serverErrChan := make(chan error, 1)

	go func() {
		s.logger.Info("Starting HTTP server",
			slog.String("port", s.cfg.Port),
			slog.String("env", s.cfg.AppEnv),
		)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- err
		}
	}()

	select {
	case err := <-serverErrChan:
		return fmt.Errorf("http server failure: %w", err)
	case <-ctx.Done():
		s.logger.Info("Shutting down server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("forced server shutdown failed: %w", err)
		}

		s.logger.Info("Server stopped cleanly")
		return nil
	}
}
