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

type Container struct {
	models.Container
}

type ContainerEnv struct {
	models.ContainerEnv

	Container *Container `gorm:"foreignKey:ContainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Volume struct {
	models.Volume
}

type ContainerVolume struct {
	models.ContainerVolume

	Container *Container `gorm:"foreignKey:ContainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Volume    *Volume    `gorm:"foreignKey:VolumeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ContainerPort struct {
	models.ContainerPort

	Container *Container `gorm:"foreignKey:ContainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DockerTask struct {
	models.DockerTask

	Container *Container `gorm:"foreignKey:ContainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func Migrate(
	ctx context.Context,
	db *gorm.DB,
) error {
	migrator := db.WithContext(ctx).Migrator()

	err := migrator.AutoMigrate(
		&User{},
		&Settings{},
		&Container{},
		&Volume{},
		&ContainerEnv{},
		&ContainerVolume{},
		&ContainerPort{},
		&DockerTask{},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
