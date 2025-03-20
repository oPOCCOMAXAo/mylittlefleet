package migrations

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Tables with relations should be defined here and used in the Migrate function.
//
// Automigrator will fail if will be multiple definitions of the same table. E.g.:
// - with relations in this package
// - without relations in the models package

type User struct {
	models.User
}

type Settings struct {
	models.Settings
}

func Migrate(
	ctx context.Context,
	db *gorm.DB,
) error {
	migrator := db.WithContext(ctx).Migrator()

	err := migrator.AutoMigrate(
		&User{},
		&Settings{},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
