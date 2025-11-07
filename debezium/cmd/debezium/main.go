package main

import (
	"context"
	"debezium_server/internal/config"
	v1 "debezium_server/internal/transport/http/v1"
	"debezium_server/pkg/logger"
	"debezium_server/pkg/postgres"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./config/.env"
	}

	if err := godotenv.Load(envPath); err != nil {
		panic(fmt.Errorf("error loading .env file: %s", err))
	}

	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	db, err := postgres.New(cfg.Config)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	lg := logger.NewLogger(cfg.Environment)

	ctx := logger.WithRequestID(context.Background(), "12345678")

	lg.Info(ctx, "starting server")

	server := v1.NewServer(cfg.Port, db.Pool)
	err = server.RegisterHandlers()
	if err != nil {
		lg.Error(ctx, "failed to register handlers", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		lg.Info(ctx, fmt.Sprintf("HTTP server listening on port %d", cfg.Port))
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			lg.Error(ctx, "server error", zap.Error(err))
		}
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	lg.Info(ctx, "Shutdown signal received, starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		lg.Error(ctx, "server shutdown error", zap.Error(err))
	}

	db.Close()
	lg.Info(ctx, "Database connection pool closed")

	wg.Wait()
	lg.Info(ctx, "Server stopped gracefully")
}
