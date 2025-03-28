package container

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
)

func (s *Service) GetContainersForDashboard(
	ctx context.Context,
) ([]*models.FullContainerInfo, error) {
	containers, err := s.repo.GetDashboardContainers(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*models.FullContainerInfo, len(containers))
	for i, container := range containers {
		res[i] = &models.FullContainerInfo{
			Container: container,
		}
	}

	return res, nil
}
