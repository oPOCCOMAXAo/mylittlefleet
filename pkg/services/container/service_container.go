package container

import (
	"context"
	"strconv"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
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

//nolint:cyclop
func (s *Service) SaveFullContainerSettings(
	ctx context.Context,
	container *models.FullContainerInfo,
) error {
	var (
		existing *models.Container
		err      error
	)

	if container.Container.ID != 0 {
		existing, err = s.repo.GetContainerByID(ctx, container.Container.ID)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return err
		}

		container.Container.ID = 0
	}

	if existing == nil {
		existing, err = s.GetContainerByName(ctx, container.Container.Name)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return err
		}
	}

	if existing == nil {
		err = s.repo.CreateContainer(ctx, container.Container)
		if err != nil {
			return err
		}
	} else {
		s.setContainerDiff(container.Container, existing)
	}

	if container.Container.DockerName == "" {
		container.Container.DockerName = s.generateContainerName(container.Container)
	}

	err = s.repo.UpdateContainer(ctx, container.Container)
	if err != nil {
		return err
	}

	err = s.SaveContainerVolumes(ctx, container.Container.ID, container.Volumes)
	if err != nil {
		return err
	}

	err = s.SaveContainerPorts(ctx, container.Container.ID, container.Ports)
	if err != nil {
		return err
	}

	err = s.SaveContainerEnvs(ctx, container.Container.ID, container.Envs)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) generateContainerName(
	container *models.Container,
) string {
	return ContainerNamePrefix +
		s.GetInstallationID() + "-" +
		strconv.FormatInt(container.ID, 10)
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

func (s *Service) FillContainers(
	ctx context.Context,
	containers ...*models.FullContainerInfo,
) error {
	for _, container := range containers {
		if container.Container.DockerName == "" {
			container.Container.DockerName = s.generateContainerName(container.Container)
		}
	}

	err := s.FillVolumesInfo(ctx, containers)
	if err != nil {
		return err
	}

	err = s.FillPortsInfo(ctx, containers)
	if err != nil {
		return err
	}

	err = s.FillEnvsInfo(ctx, containers)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) getAllFullContainerInfos(
	ctx context.Context,
) ([]*models.FullContainerInfo, error) {
	containers, err := s.repo.GetAllContainers(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*models.FullContainerInfo, 0, len(containers))
	for _, container := range containers {
		res = append(res, &models.FullContainerInfo{
			Container: container,
		})
	}

	err = s.FillContainers(ctx, res...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetContainerFullInfoByID(
	ctx context.Context,
	id int64,
) (*models.FullContainerInfo, error) {
	container, err := s.repo.GetContainerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &models.FullContainerInfo{
		Container: container,
	}

	err = s.FillContainers(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) StartContainerByName(
	ctx context.Context,
	name string,
) error {
	containerID, err := s.repo.GetContainerIDByName(ctx, name)
	if err != nil {
		return err
	}

	return s.StartContainerByID(ctx, containerID)
}

func (s *Service) StopContainerByName(
	ctx context.Context,
	name string,
) error {
	containerID, err := s.repo.GetContainerIDByName(ctx, name)
	if err != nil {
		return err
	}

	return s.StopContainerByID(ctx, containerID)
}

func (s *Service) StartContainerByID(
	ctx context.Context,
	id int64,
) error {
	cnt, err := s.repo.GetContainerByID(ctx, id)
	if err != nil {
		return err
	}

	if cnt.Paused {
		err = s.repo.UpdateContainerPaused(ctx, id, false)
		if err != nil {
			return err
		}
	}

	return s.CreateTasks(ctx, &models.DockerTask{
		ContainerID: id,
		Action:      models.DTAStart,
	})
}

func (s *Service) StopContainerByID(
	ctx context.Context,
	id int64,
) error {
	cnt, err := s.repo.GetContainerByID(ctx, id)
	if err != nil {
		return err
	}

	if !cnt.Paused {
		err = s.repo.UpdateContainerPaused(ctx, id, true)
		if err != nil {
			return err
		}
	}

	return s.CreateTasks(ctx, &models.DockerTask{
		ContainerID: id,
		Action:      models.DTAStop,
	})
}
