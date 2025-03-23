package container

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/structs"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

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

func (s *Service) EnsureFullContainerSettings(
	ctx context.Context,
	params *structs.FullContainerInfo,
) error {
	exiting, err := s.GetContainerByName(ctx, params.Container.Name)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	if exiting == nil {
		err = s.repo.CreateContainer(ctx, params.Container)
		if err != nil {
			return err
		}
	} else {
		s.setContainerDiff(params.Container, exiting)

		err = s.repo.UpdateContainer(ctx, params.Container)
		if err != nil {
			return err
		}
	}

	err = s.EnsureContainerVolumes(ctx, params.Container.ID, params.Volumes)
	if err != nil {
		return err
	}

	err = s.EnsureContainerPorts(ctx, params.Container.ID, params.Ports)
	if err != nil {
		return err
	}

	err = s.EnsureContainerEnvs(ctx, params.Container.ID, params.Envs)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) setContainerDiff(newItem, oldItem *models.Container) {
	newItem.ID = oldItem.ID
	newItem.DockerID = lo.CoalesceOrEmpty(newItem.DockerID, oldItem.DockerID)
	newItem.DockerName = lo.CoalesceOrEmpty(newItem.DockerName, oldItem.DockerName)
	newItem.Image = lo.CoalesceOrEmpty(newItem.Image, oldItem.Image)
	newItem.Tag = lo.CoalesceOrEmpty(newItem.Tag, oldItem.Tag)
	newItem.Paused = lo.CoalesceOrEmpty(newItem.Paused, oldItem.Paused)
	newItem.Deleted = lo.CoalesceOrEmpty(newItem.Deleted, oldItem.Deleted)
	newItem.Internal = lo.CoalesceOrEmpty(newItem.Internal, oldItem.Internal)
}
