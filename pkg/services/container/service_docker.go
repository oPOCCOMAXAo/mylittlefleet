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

func (s *Service) serveSyncWithDocker() {
	s.RunSyncWithDocker()

	for {
		select {
		case <-s.runCtx.Done():
			return
		case <-s.chanSyncWithDocker:
			err := s.syncWithDocker(s.runCtx)
			if err != nil {
				s.logger.Error("syncWithDocker", slog.Any("error", err))
			}
		}
	}
}

func (s *Service) RunSyncWithDocker() {
	select {
	case s.chanSyncWithDocker <- struct{}{}:
	default:
	}
}

func (s *Service) syncWithDocker(ctx context.Context) error {
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

// dockerContainer may be nil when container is not running.
//
//nolint:mnd
func (s *Service) syncCurrentInstallation(
	ctx context.Context,
	dockerContainer *container.Summary,
) error {
	oldContainer, err := s.GetContainerByName(ctx, ContainerNameSelf)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	if oldContainer == nil {
		return errors.WithStack(models.ErrFlowBroken)
	}

	if dockerContainer == nil {
		oldContainer.Paused = true
	} else {
		oldContainer.DockerID = dockerContainer.ID
		oldContainer.DockerName = lo.CoalesceOrEmpty(dockerContainer.Names...)

		image := strings.Split(dockerContainer.Image, ":")
		if len(image) >= 2 {
			oldContainer.Image = image[0]
			oldContainer.Tag = image[1]
		} else {
			oldContainer.Image = dockerContainer.Image
		}
	}

	return s.repo.UpdateContainer(ctx, oldContainer)
}
