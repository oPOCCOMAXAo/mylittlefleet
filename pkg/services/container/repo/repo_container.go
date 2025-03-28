package repo

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) CreateContainer(
	ctx context.Context,
	container *models.Container,
) error {
	err := r.db.
		WithContext(ctx).
		Create(container).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetContainersByName(
	ctx context.Context,
	names ...string,
) ([]*models.Container, error) {
	if len(names) == 0 {
		return nil, nil
	}

	var res []*models.Container

	err := r.db.
		WithContext(ctx).
		Where("name IN (?)", names).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (r *Repo) GetContainerByID(
	ctx context.Context,
	id int64,
) (*models.Container, error) {
	var res models.Container

	err := r.db.
		WithContext(ctx).
		Model(&models.Container{}).
		Where("id = ?", id).
		Take(&res).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(models.ErrNotFound)
		}

		return nil, errors.WithStack(err)
	}

	return &res, nil
}

func (r *Repo) UpdateContainer(
	ctx context.Context,
	value *models.Container,
) error {
	err := r.db.
		WithContext(ctx).
		Model(&models.Container{}).
		Where("id = ?", value.ID).
		Select("*").
		Updates(value).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetAllContainers(
	ctx context.Context,
) ([]*models.Container, error) {
	var res []*models.Container

	err := r.db.
		WithContext(ctx).
		Model(&models.Container{}).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (r *Repo) UpdateContainerDockerID(
	ctx context.Context,
	id int64,
	dockerID string,
) error {
	err := r.db.
		WithContext(ctx).
		Model(&models.Container{}).
		Where("id = ?", id).
		Update("docker_id", dockerID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetContainerIDByName(
	ctx context.Context,
	name string,
) (int64, error) {
	var res models.Container

	err := r.db.
		WithContext(ctx).
		Model(&models.Container{}).
		Where("name = ?", name).
		Select("id").
		Take(&res).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.WithStack(models.ErrNotFound)
		}

		return 0, errors.WithStack(err)
	}

	return res.ID, nil
}

func (r *Repo) UpdateContainerPaused(
	ctx context.Context,
	id int64,
	paused bool,
) error {
	err := r.db.
		WithContext(ctx).
		Model(&models.Container{}).
		Where("id = ?", id).
		Update("paused", paused).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
