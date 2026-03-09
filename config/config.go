package config

import (
	"log/slog"
	"online-learning-platform-go-api/internal/pkg"
	"os"

	"github.com/joho/godotenv"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Minio    pkg.MinioConfig    `yaml:"minio"`
	Postgres pkg.PostgresConfig `yaml:"postgres"`
	Server   pkg.Server         `yaml:"server"`
	Token    pkg.TokenConfig    `yaml:"token"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load("../.env")

	data, err := os.ReadFile("../resources/config.yaml")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	expanded := os.ExpandEnv(string(data))

	cfg := &Config{}

	if err := yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return cfg, nil
}
