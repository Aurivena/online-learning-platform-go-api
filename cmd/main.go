package main

import (
	"online-learning-platform-go-api/internal/pkg"

	"online-learning-platform-go-api/config"
	"online-learning-platform-go-api/internal/gateway"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		return
	}

	db, err := pkg.NewPostgresConfig(cfg.Postgres)
	if err != nil {
		return
	}

	defer db.Close()

	_, err = pkg.NewMinioConfig(cfg.Minio)
	if err != nil {
		return
	}

	handler := gateway.NewGateway(cfg.Server)

	go pkg.RunServer(cfg.Server, handler)

	pkg.StopServer()
}
