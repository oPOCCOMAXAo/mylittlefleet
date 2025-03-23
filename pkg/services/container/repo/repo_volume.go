package repo

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) CreateInternalContainerVolume(
	ctx context.Context,
	cVolume *models.ContainerVolume,
) error {
	err := r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			volume := models.Volume{
				Internal: true,
			}

			err := tx.Create(&volume).Error
			if err != nil {
				return errors.WithStack(err)
			}

			cVolume.VolumeID = volume.ID

			err = tx.
				Create(cVolume).
				Error
			if err != nil {
				return errors.WithStack(err)
			}

			return nil
		})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetContainerVolumes(
	ctx context.Context,
	containerID int64,
) ([]*models.ContainerVolume, error) {
	var res []*models.ContainerVolume

	err := r.db.
		WithContext(ctx).
		Model(&models.ContainerVolume{}).
		Where("container_id = ?", containerID).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (r *Repo) DeleteContainerVolumesByID(
	ctx context.Context,
	ids []int64,
) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Where("id IN (?)", ids).
		Delete(&models.ContainerVolume{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) UpdateContainerVolumes(
	ctx context.Context,
	volumes []*models.ContainerVolume,
) error {
	if len(volumes) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			for _, volume := range volumes {
				err := tx.
					Model(&models.ContainerVolume{}).
					Where("id = ?", volume.ID).
					Select("*").
					Updates(volume).
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
