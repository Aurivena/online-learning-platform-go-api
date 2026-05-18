package pkg

import (
	"database/sql"
	"fmt"
	"log/slog"
	"online-learning-platform-go-api/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host, Port, User, Password, DBName, SSL string
}

func NewPostgresConfig(cfg DBConfig) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSL)

	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("open postgres connection", slog.String("error", err.Error()))
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		slog.Error("ping postgres", slog.String("error", err.Error()))
		return nil, err
	}

	if err := goose.Up(sqlDB, config.MigrationsDirForPostgres); err != nil {

		slog.Error(err.Error())
		return nil, err
	}

	return gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

}
