package repo

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) GetContainerPorts(
	ctx context.Context,
	containerID int64,
) ([]*models.ContainerPort, error) {
	var res []*models.ContainerPort

	err := r.db.
		WithContext(ctx).
		Model(&models.ContainerPort{}).
		Where("container_id = ?", containerID).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (r *Repo) CreateContainerPorts(
	ctx context.Context,
	ports []*models.ContainerPort,
) error {
	if len(ports) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Create(ports).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) DeleteContainerPortsByID(
	ctx context.Context,
	ids []int64,
) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Where("id IN (?)", ids).
		Delete(&models.ContainerPort{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) UpdateContainerPorts(
	ctx context.Context,
	ports []*models.ContainerPort,
) error {
	if len(ports) == 0 {
		return nil
	}

	err := r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			for _, port := range ports {
				err := tx.
					Model(&models.ContainerPort{}).
					Where("id = ?", port.ID).
					Select("*").
					Updates(port).
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
