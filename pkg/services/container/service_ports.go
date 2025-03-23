package container

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/diff"
)

func (s *Service) EnsureContainerPorts(
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
