package migrations

import (
	"context"

	"gorm.io/gorm"
)

func Migrate(
	ctx context.Context,
	db *gorm.DB,
) error {
	return nil
}
