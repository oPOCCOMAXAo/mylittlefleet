package container

import (
	"context"
	"strconv"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/diff"
)

//nolint:cyclop,funlen
func (s *Service) SaveContainerVolumes(
	ctx context.Context,
	containerID int64,
	domain []*models.VolumeDomain,
) error {
	err := s.SaveVolumes(ctx, domain)
	if err != nil {
		return err
	}

	oldVolumes, err := s.repo.GetContainerVolumes(ctx, containerID)
	if err != nil {
		return err
	}

	names := make([]string, 0, len(domain))

	for _, v := range domain {
		if v.Volume.DockerName != "" {
			names = append(names, v.Volume.DockerName)
		}
	}

	rawVolumes, err := s.repo.GetVolumesByDockerName(ctx, names)
	if err != nil {
		return err
	}

	rawIDByName := make(map[string]int64, len(rawVolumes))
	for _, v := range rawVolumes {
		rawIDByName[v.DockerName] = v.ID
	}

	newVolumes := make([]*models.ContainerVolume, 0, len(domain))
	for _, v := range domain {
		newVolumes = append(newVolumes, &models.ContainerVolume{
			ContainerID:   containerID,
			VolumeID:      rawIDByName[v.Volume.DockerName],
			ContainerPath: v.ContainerPath,
		})
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

func (s *Service) SaveVolumes(
	ctx context.Context,
	domain []*models.VolumeDomain,
) error {
	ids := make([]int64, 0, len(domain))

	for _, v := range domain {
		if v.Volume.ID != 0 {
			ids = append(ids, v.Volume.ID)
		}
	}

	oldVolumes, err := s.repo.GetVolumesByID(ctx, ids)
	if err != nil {
		return err
	}

	newVolumes := make([]*models.Volume, 0, len(domain))

	for _, v := range domain {
		if v.Volume.ID != 0 && v.Volume.DockerName != "" {
			newVolumes = append(newVolumes, &v.Volume)
		}
	}

	diff := diff.Slices(
		newVolumes,
		oldVolumes,
		(*models.Volume).UniqueKey,
		(*models.Volume).Equal,
		(*models.Volume).PrepareForUpdate,
	)

	err = s.repo.UpdateVolumes(ctx, diff.Updated)
	if err != nil {
		return err
	}

	// delete will be done in cleanup procedure.
	// create will be done in the next step.

	return nil
}

func (s *Service) FillVolumesInfo(
	ctx context.Context,
	containers []*models.FullContainerInfo,
) error {
	ids := make([]int64, 0, len(containers))
	for _, container := range containers {
		ids = append(ids, container.Container.ID)
	}

	cVol, err := s.repo.GetContainerVolumesByContainerIDs(ctx, ids)
	if err != nil {
		return err
	}

	volIDs := make([]int64, 0, len(cVol))

	cvByContainerID := make(map[int64][]*models.ContainerVolume, len(containers))
	for _, v := range cVol {
		cvByContainerID[v.ContainerID] = append(cvByContainerID[v.ContainerID], v)
		volIDs = append(volIDs, v.VolumeID)
	}

	volumes, err := s.repo.GetVolumesByID(ctx, volIDs)
	if err != nil {
		return err
	}

	volByID := make(map[int64]*models.Volume, len(volumes))

	for _, v := range volumes {
		if v.DockerName == "" {
			v.DockerName = s.generateVolumeName(v)
		}

		volByID[v.ID] = v
	}

	for _, container := range containers {
		cvs := cvByContainerID[container.Container.ID]
		container.Volumes = make([]*models.VolumeDomain, 0, len(cvs))

		for _, cv := range cvs {
			temp := &models.VolumeDomain{
				ContainerVolume: *cv,
			}

			v := volByID[cv.VolumeID]
			if v != nil {
				temp.Volume = *v
			}

			container.Volumes = append(container.Volumes, temp)
		}
	}

	return nil
}

func (s *Service) generateVolumeName(v *models.Volume) string {
	return VolumeNamePrefix +
		s.GetInstallationID() + "-" +
		strconv.FormatInt(v.ID, 10)
}
