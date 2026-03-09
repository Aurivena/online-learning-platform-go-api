package pkg

import (
	"database/sql"
	"fmt"
	"log/slog"
	"online-learning-platform-go-api/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func NewPostgresConfig(cfg config.PostgresConfig) (*sqlx.DB, error) {
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

	if err := goose.Up(db, config.MigrationsDirForPostgres); err != nil {

		slog.Error(err.Error())
		return nil, err
	}

	return sqlx.NewDb(db, "postgres"), nil
}
