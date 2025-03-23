package container

import (
	"context"
	"slices"
	"strings"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/diff"
)

func (s *Service) EnsureContainerEnvs(
	ctx context.Context,
	containerID int64,
	newEnvs []*models.ContainerEnv,
) error {
	oldEnvs, err := s.repo.GetContainerEnvs(ctx, containerID)
	if err != nil {
		return err
	}

	for _, env := range newEnvs {
		env.ContainerID = containerID
	}

	slices.SortFunc(newEnvs, func(l, r *models.ContainerEnv) int {
		return strings.Compare(l.Name, r.Name)
	})

	diff := diff.Slices(
		newEnvs,
		oldEnvs,
		(*models.ContainerEnv).UniqueKey,
		(*models.ContainerEnv).Equal,
		(*models.ContainerEnv).PrepareForUpdate,
	)

	err = s.repo.CreateContainerEnvs(ctx, diff.Created)
	if err != nil {
		return err
	}

	ids := make([]int64, 0, len(diff.Deleted))
	for _, env := range diff.Deleted {
		ids = append(ids, env.ID)
	}

	err = s.repo.DeleteContainerEnvsByID(ctx, ids)
	if err != nil {
		return err
	}

	err = s.repo.UpdateContainerEnvs(ctx, diff.Updated)
	if err != nil {
		return err
	}

	return nil
}
