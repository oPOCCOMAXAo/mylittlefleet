package logger

import (
	"log/slog"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"
)

func DecorateGormDB(
	logger *slog.Logger,
	db *gorm.DB,
) *gorm.DB {
	return db.Session(&gorm.Session{
		Logger: slogGorm.New(
			slogGorm.WithHandler(logger.Handler()),
		),
	})
}
