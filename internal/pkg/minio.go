package pkg

import (
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	AccessKey, SecretKey, Endpoint string
	SSL                            bool
}

func NewMinioConfig(cfg MinioConfig) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		slog.Error("Failed to create Minio client: ", "error", err)
		return nil, err
	}

	return minioClient, nil
}
