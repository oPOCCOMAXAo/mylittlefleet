package db

import (
	"log/slog"

	slogGorm "github.com/orandin/slog-gorm"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DSN string `env:"DSN,required"` // sqlite path: "./db.sqlite3?params"
}

func NewSQLite(
	cfg Config,
	logger *slog.Logger,
) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithHandler(logger.Handler()),
		),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
