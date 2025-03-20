package repo

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetSettingsByKeys(
	ctx context.Context,
	keys ...string,
) ([]*models.Settings, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	var res []*models.Settings

	err := r.db.WithContext(ctx).
		Where("key IN ?", keys).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (r *Repo) UpdateSettings(
	ctx context.Context,
	settings ...*models.Settings,
) error {
	if len(settings) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Save(&settings).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) DeleteSettingsByKeys(
	ctx context.Context,
	keys ...string,
) error {
	if len(keys) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Where("key IN ?", keys).
		Delete(&models.Settings{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
