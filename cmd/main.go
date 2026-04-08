package main

import (
	"log/slog"
	"online-learning-platform-go-api/internal/di"
	"online-learning-platform-go-api/internal/middleware"
	"online-learning-platform-go-api/internal/pkg"

	"online-learning-platform-go-api/config"
	"online-learning-platform-go-api/internal/gateway"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		return
	}

	gorm, err := pkg.NewPostgresConfig(pkg.DBConfig{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DBName:   cfg.Postgres.Database,
		SSL:      cfg.Postgres.SSL,
	})
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL: ", "error", err)
		return
	}

	sqlDB, err := gorm.DB()
	if err != nil {
		slog.Error("Failed to get SQL DB: ", "error", err)
		return
	}

	defer sqlDB.Close()

	_, err = pkg.NewMinioConfig(pkg.MinioConfig{
		AccessKey: cfg.Minio.AccessKey,
		SecretKey: cfg.Minio.SecretKey,
		Endpoint:  cfg.Minio.Endpoint,
		SSL:       cfg.Minio.SSL,
	})
	if err != nil {
		return
	}

	provider := di.NewProvider(gorm)

	router := gateway.SetupRouter(cfg.Server, middleware.NewMiddleware(&cfg.Token), gateway.NewGateway(provider))

	httpServer := pkg.RunServer(cfg.Server.Addr, cfg.Server.Port, router)

	pkg.StopServer(httpServer)
}
