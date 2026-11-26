package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ticketflow/ticketflow/internal/config"
	"github.com/ticketflow/ticketflow/internal/infrastructure/postgres"
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

	// 4. Initialize PostgreSQL connection pool
	dbPool, err := postgres.NewPool(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbPool.Close()

	logger.Info("Connected to PostgreSQL",
		slog.String("host", cfg.DB.Host),
		slog.String("port", cfg.DB.Port),
		slog.String("database", cfg.DB.DBName),
		slog.String("schema", cfg.DB.Schema),
	)

	// 5. Initialize and run HTTP server
	srv := server.NewServer(cfg, logger, dbPool)
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
