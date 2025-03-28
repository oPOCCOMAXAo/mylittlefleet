package container

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/diff"
)

func (s *Service) SaveContainerPorts(
	ctx context.Context,
	containerID int64,
	newPorts []*models.ContainerPort,
) error {
	oldPorts, err := s.repo.GetContainerPorts(ctx, containerID)
	if err != nil {
		return err
	}

	for _, port := range newPorts {
		port.ContainerID = containerID
	}

	diff := diff.Slices(
		newPorts,
		oldPorts,
		(*models.ContainerPort).UniqueKey,
		(*models.ContainerPort).Equal,
		(*models.ContainerPort).PrepareForUpdate,
	)

	err = s.repo.CreateContainerPorts(ctx, diff.Created)
	if err != nil {
		return err
	}

	ids := make([]int64, 0, len(diff.Deleted))
	for _, port := range diff.Deleted {
		ids = append(ids, port.ID)
	}

	err = s.repo.DeleteContainerPortsByID(ctx, ids)
	if err != nil {
		return err
	}

	err = s.repo.UpdateContainerPorts(ctx, diff.Updated)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FillPortsInfo(
	ctx context.Context,
	containers []*models.FullContainerInfo,
) error {
	ids := make([]int64, 0, len(containers))
	for _, container := range containers {
		ids = append(ids, container.Container.ID)
	}

	ports, err := s.repo.GetContainerPortsByContainerIDs(ctx, ids)
	if err != nil {
		return err
	}

	portsByContainerID := make(map[int64][]*models.ContainerPort, len(ports))
	for _, port := range ports {
		portsByContainerID[port.ContainerID] = append(portsByContainerID[port.ContainerID], port)
	}

	for _, container := range containers {
		container.Ports = portsByContainerID[container.Container.ID]
	}

	return nil
}
