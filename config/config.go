package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"go.yaml.in/yaml/v3"
)

var MigrationsDirForPostgres = "../resources/migrations"

type Config struct {
	Minio    MinioConfig    `yaml:"minio"`
	Postgres PostgresConfig `yaml:"postgres"`
	Server   Server         `yaml:"server"`
	Token    TokenConfig    `yaml:"token"`
}

type MinioConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Endpoint  string `yaml:"endpoint"`
	SSL       bool   `yaml:"sslmode"`
}

type PostgresConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Port     string `yaml:"port"`
	SSL      string `yaml:"sslmode"`
}

type Server struct {
	Addr string
	Port string
}

type TokenConfig struct {
	AccessToken  string `yaml:"access-token"`
	RefreshToken string `yaml:"refresh-token"`
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
