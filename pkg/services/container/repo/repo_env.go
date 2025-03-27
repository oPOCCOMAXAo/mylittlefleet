package repo

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) GetContainerEnvs(
	ctx context.Context,
	containerID int64,
) ([]*models.ContainerEnv, error) {
	var res []*models.ContainerEnv

	err := r.db.
		WithContext(ctx).
		Model(&models.ContainerEnv{}).
		Where("container_id = ?", containerID).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (r *Repo) CreateContainerEnvs(
	ctx context.Context,
	envs []*models.ContainerEnv,
) error {
	if len(envs) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Create(envs).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) DeleteContainerEnvsByID(
	ctx context.Context,
	ids []int64,
) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Where("id IN (?)", ids).
		Delete(&models.ContainerEnv{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) UpdateContainerEnvs(
	ctx context.Context,
	envs []*models.ContainerEnv,
) error {
	if len(envs) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			for _, env := range envs {
				err := tx.
					Model(&models.ContainerEnv{}).
					Where("id = ?", env.ID).
					Select("*").
					Updates(env).
					Error
				if err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetContainerEnvsByContainerIDs(
	ctx context.Context,
	ids []int64,
) ([]*models.ContainerEnv, error) {
	var res []*models.ContainerEnv

	err := r.db.
		WithContext(ctx).
		Model(&models.ContainerEnv{}).
		Where("container_id IN (?)", ids).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}
