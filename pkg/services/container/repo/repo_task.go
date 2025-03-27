package repo

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) CreateTasks(
	ctx context.Context,
	tasks []*models.DockerTask,
) error {
	err := r.db.
		WithContext(ctx).
		Create(tasks).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetFirstTaskToExecute(
	ctx context.Context,
) (*models.DockerTask, error) {
	var res models.DockerTask

	err := r.db.
		WithContext(ctx).
		Model(&models.DockerTask{}).
		Where("finished = 0").
		Order("id ASC").
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

func (r *Repo) UpdateTask(
	ctx context.Context,
	task *models.DockerTask,
) error {
	err := r.db.
		WithContext(ctx).
		Model(&models.DockerTask{}).
		Where("id = ?", task.ID).
		Select("*").
		Updates(task).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
