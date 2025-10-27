package main

import (
	"context"
	"debezium_server/internal/config"
	"debezium_server/pkg/logger"
)

func main() {
	//envPath := flag.String("env", ".env", "path to the environment file")
	//flag.Parse()
	cfg, err := config.ParseConfig("../../config/.env")
	if err != nil {
		return
	}
	lg := logger.NewLogger(cfg.Environment)
	ctx := logger.WithRequestID(context.Background(), "12345678")

	lg.Info(ctx, "starting server")
}
