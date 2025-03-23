package container

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/diff"
)

func (s *Service) EnsureContainerVolumes(
	ctx context.Context,
	containerID int64,
	newVolumes []*models.ContainerVolume,
) error {
	oldVolumes, err := s.repo.GetContainerVolumes(ctx, containerID)
	if err != nil {
		return err
	}

	for _, volume := range newVolumes {
		volume.ContainerID = containerID
	}

	diff := diff.Slices(
		newVolumes,
		oldVolumes,
		(*models.ContainerVolume).UniqueKey,
		(*models.ContainerVolume).Equal,
		(*models.ContainerVolume).PrepareForUpdate,
	)

	for _, cv := range diff.Created {
		err = s.repo.CreateInternalContainerVolume(ctx, cv)
		if err != nil {
			return err
		}
	}

	ids := make([]int64, 0, len(diff.Deleted))
	for _, cv := range diff.Deleted {
		ids = append(ids, cv.ID)
	}

	err = s.repo.DeleteContainerVolumesByID(ctx, ids)
	if err != nil {
		return err
	}

	err = s.repo.UpdateContainerVolumes(ctx, diff.Updated)
	if err != nil {
		return err
	}

	return nil
}
