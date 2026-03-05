package pkg

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var migrationsDir = "../resources/migrations"

type PostgresConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Port     string `yaml:"port"`
	SSL      string `yaml:"sslmode"`
}

func NewPostgresConfig(cfg PostgresConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSL)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("open postgres connection", slog.String("error", err.Error()))
		return nil, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("ping postgres", slog.String("error", err.Error()))
		return nil, err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return sqlx.NewDb(db, "postgres"), nil
}
