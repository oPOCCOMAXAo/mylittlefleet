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
	return &Repo{db: db}
}

func (r *Repo) GetTotalUsers(
	ctx context.Context,
) (int64, error) {
	var res int64

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Count(&res).
		Error
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return res, nil
}

func (r *Repo) CreateUser(
	ctx context.Context,
	user *models.User,
) error {
	err := r.db.WithContext(ctx).
		Create(user).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetUserByID(
	ctx context.Context,
	id int64,
) (*models.User, error) {
	var res models.User

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
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

func (r *Repo) GetUserByLogin(
	ctx context.Context,
	login string,
) (*models.User, error) {
	var res models.User

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("login = ?", login).
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
