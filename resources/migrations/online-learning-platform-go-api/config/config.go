package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.yaml.in/yaml/v3"
)

var (
	MigrationsDirForPostgres = "resources/migrations"
	AccessTokenTimestamp     = 0
	RefreshTokenTimestamp    = 0
)

type Config struct {
	Minio    MinioConfig    `yaml:"minio"`
	Postgres PostgresConfig `yaml:"postgres"`
	Server   Server         `yaml:"server"`
	Token    TokenConfig    `yaml:"token"`
}

type MinioConfig struct {
	AccessKey     string `yaml:"access_key"`
	SecretKey     string `yaml:"secret_key"`
	Bucket        string `yaml:"bucket"`
	Endpoint      string `yaml:"endpoint"`
	SSL           bool   `yaml:"sslmode"`
	PublicBaseURL string `yaml:"public_base_url"` // optional, e.g. http://127.0.0.1:9000 — for payload.url (path-style /bucket/key)
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
	Addr         string
	Port         string
	ServerDomain string `yaml:"server-domain"`
}

type TokenConfig struct {
	AccessToken  string `yaml:"access-token"`
	RefreshToken string `yaml:"refresh-token"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load(findEnvFile())

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		var err error
		configPath, err = findProjectPath("resources/config.yaml")
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
	}

	migrationsDir, err := findProjectPath("resources/migrations")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	MigrationsDirForPostgres = migrationsDir

	data, err := os.ReadFile(configPath)
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

func findProjectPath(rel string) (string, error) {
	if filepath.IsAbs(rel) {
		return rel, nil
	}
	dir, _ := os.Getwd()
	for {
		path := filepath.Join(dir, rel)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("project path %q not found from current directory", rel)
}

func findEnvFile() string {
	dir, _ := os.Getwd()
	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ".env"
}
