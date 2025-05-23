package container

import (
	"context"
	"slices"
	"strings"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/structs"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/diff"
)

func (s *Service) SaveContainerEnvs(
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

func (s *Service) FillEnvsInfo(
	ctx context.Context,
	containers []*structs.FullContainerInfo,
) error {
	ids := make([]int64, 0, len(containers))
	for _, container := range containers {
		ids = append(ids, container.Container.ID)
	}

	envs, err := s.repo.GetContainerEnvsByContainerIDs(ctx, ids)
	if err != nil {
		return err
	}

	envsByContainerID := make(map[int64][]*models.ContainerEnv)
	for _, env := range envs {
		envsByContainerID[env.ContainerID] = append(envsByContainerID[env.ContainerID], env)
	}

	for _, container := range containers {
		container.Envs = envsByContainerID[container.Container.ID]
	}

	return nil
}
