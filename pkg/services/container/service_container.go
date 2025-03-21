package container

import (
	"context"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (s *Service) serveSyncContainers() {
	s.TriggerSyncContainers()

	for {
		select {
		case <-s.runCtx.Done():
			return
		case <-s.chanSyncContainers:
			err := s.syncContainers(s.runCtx)
			if err != nil {
				s.logger.Error("syncContainers", slog.Any("error", err))
			}
		}
	}
}

func (s *Service) syncContainers(ctx context.Context) error {
	current, err := s.FindCurrentInstallation(ctx)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	err = s.syncCurrentInstallation(ctx, current)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) syncCurrentInstallation(
	ctx context.Context,
	docker *container.Summary,
) error {
	db, err := s.GetContainerByName(ctx, ContainerNameSelf)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	if db == nil {
		return errors.WithStack(models.ErrFlowBroken)
	}

	if docker == nil {
		db.Paused = true
	} else {
		db.DockerID = docker.ID
		db.DockerName = lo.CoalesceOrEmpty(docker.Names...)

		image := strings.Split(docker.Image, ":")
		if len(image) >= 2 {
			db.Image = image[0]
			db.Tag = image[1]
		} else {
			db.Image = docker.Image
		}
	}

	return s.repo.UpdateContainer(ctx, db)
}

func (s *Service) TriggerSyncContainers() {
	select {
	case s.chanSyncContainers <- struct{}{}:
	default:
	}
}

func (s *Service) GetContainersMapByName(
	ctx context.Context,
	names ...string,
) (map[string]*models.Container, error) {
	containers, err := s.repo.GetContainersByName(ctx, names...)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*models.Container, len(containers))
	for _, container := range containers {
		res[container.Name] = container
	}

	return res, nil
}

func (s *Service) GetContainerByName(
	ctx context.Context,
	name string,
) (*models.Container, error) {
	containers, err := s.repo.GetContainersByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, errors.WithStack(models.ErrNotFound)
	}

	return containers[0], nil
}
