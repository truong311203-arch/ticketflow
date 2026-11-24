package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ticketflow/ticketflow/internal/config"
	"github.com/ticketflow/ticketflow/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "application startup error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// 2. Initialize structured logger
	logger := initLogger(cfg)
	slog.SetDefault(logger)

	logger.Info("Starting TicketFlow API",
		slog.String("environment", cfg.AppEnv),
		slog.String("port", cfg.Port),
	)

	// 3. Setup signal trap for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. Initialize and run HTTP server
	srv := server.NewServer(cfg, logger)
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func initLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}

	if cfg.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLogLevel(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
