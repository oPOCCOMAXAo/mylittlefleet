package db

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DSN string `env:"DSN,required"`
}

func NewSQLite(
	cfg Config,
) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
