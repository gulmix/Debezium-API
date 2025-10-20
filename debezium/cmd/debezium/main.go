package main

import (
	"debezium_server/internal/config"
	"debezium_server/internal/repository"
	"debezium_server/internal/service"
	v1 "debezium_server/internal/transport/http"
	"flag"
)

func main() {
	envPath := flag.String("env", ".env", "path to the environment file")
	flag.Parse()
	cfg, err := config.ParseConfig(*envPath)
	if err != nil {
		return
	}
	up := repository.NewUserRepository()
	us := service.NewUserService(up)
	server := v1.NewServer(cfg.Port)
	if err := server.Start(); err != nil {
		return
	}
}
