package repo

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
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
