package main

import (
	"context"
	"log/slog"
	"online-learning-platform-go-api/internal/di"
	"online-learning-platform-go-api/internal/middleware"
	"online-learning-platform-go-api/internal/pkg"
	"online-learning-platform-go-api/internal/storage"

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

	minioClient, err := pkg.NewMinioConfig(pkg.MinioConfig{
		AccessKey: cfg.Minio.AccessKey,
		SecretKey: cfg.Minio.SecretKey,
		Endpoint:  cfg.Minio.Endpoint,
		SSL:       cfg.Minio.SSL,
	})
	if err != nil {
		slog.Error("Failed to create MinIO client: ", "error", err)
		return
	}

	fileBucket := storage.NewBucket(minioClient, cfg.Minio.Bucket)
	if err := fileBucket.EnsureExists(context.Background()); err != nil {
		slog.Warn("MinIO bucket ensure failed (uploads may fail until bucket exists)", "error", err)
	}

	provider := di.NewProvider(gorm)

	userGateway := gateway.NewGateway(provider.User())
	orgGateway := provider.OrganizationGateway()
	courseGateway := gateway.NewCourseGateway(
		provider.Course(),
		provider.Module(),
		provider.Slide(),
		provider.Organization(),
		fileBucket,
		cfg.Minio.PublicBaseURL,
	)

	router := gateway.SetupRouter(cfg.Server, middleware.NewMiddleware(&cfg.Token), userGateway, orgGateway, courseGateway)

	httpServer := pkg.RunServer(cfg.Server.Addr, cfg.Server.Port, router)

	pkg.StopServer(httpServer)
}
