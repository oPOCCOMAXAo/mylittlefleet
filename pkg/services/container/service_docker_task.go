package container

import (
	"context"
	"io"
	"log/slog"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/structs"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (s *Service) serveTasks() {
	s.RunTasks()

	for {
		select {
		case <-s.runCtx.Done():
			return
		case <-s.chanTasks:
			err := s.execTasks(s.runCtx)
			if err != nil {
				s.logger.Error("execTasks", slog.Any("error", err))
			}
		}
	}
}

func (s *Service) RunTasks() {
	select {
	case s.chanTasks <- struct{}{}:
	default:
	}
}

func (s *Service) CreateTasks(
	ctx context.Context,
	tasks ...*models.DockerTask,
) error {
	if len(tasks) == 0 {
		return nil
	}

	defer s.RunTasks()

	return s.repo.CreateTasks(ctx, tasks)
}

func (s *Service) execTasks(ctx context.Context) error {
	execCount := 0

	defer func() {
		if execCount > 0 {
			s.RunSyncWithDocker()
		}
	}()

	for {
		err := s.execFirstTask(ctx)
		if err == nil {
			execCount++

			continue
		}

		if errors.Is(err, models.ErrNotFound) {
			return nil
		}

		return err
	}
}

//nolint:cyclop
func (s *Service) execFirstTask(ctx context.Context) error {
	task, err := s.repo.GetFirstTaskToExecute(ctx)
	if err != nil {
		return err
	}

	switch task.Action {
	case models.DTAStop:
		err = s.stopDockerContainer(ctx, task.ContainerID)
		if err != nil {
			return err
		}
	case models.DTAStart:
		err = s.startDockerContainer(ctx, task.ContainerID)
		if err != nil {
			return err
		}
	case models.DTADelete:
		err = s.deleteDockerContainer(ctx, task.ContainerID)
		if err != nil {
			return err
		}
	case models.DTACreate:
		err = s.createDockerContainer(ctx, task.ContainerID)
		if err != nil {
			return err
		}
	}

	task.Finished = true

	err = s.repo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) stopDockerContainer(
	ctx context.Context,
	id int64,
) error {
	cnt, err := s.repo.GetContainerByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.docker.ContainerStop(ctx, cnt.DockerID, container.StopOptions{
		Signal:  "SIGTERM",
		Timeout: lo.ToPtr(30), //nolint:mnd
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) startDockerContainer(
	ctx context.Context,
	id int64,
) error {
	cnt, err := s.repo.GetContainerByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.docker.ContainerStart(ctx, cnt.DockerID, container.StartOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) deleteDockerContainer(
	ctx context.Context,
	id int64,
) error {
	cnt, err := s.repo.GetContainerByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.docker.ContainerRemove(ctx, cnt.DockerID, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) createDockerContainer(
	ctx context.Context,
	id int64,
) error {
	exists, err := s.isContainerExist(ctx, id)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	info, err := s.GetContainerFullInfoByID(ctx, id)
	if err != nil {
		return err
	}

	for _, v := range info.Volumes {
		if v.Volume.DockerName == "" {
			err = s.createVolume(ctx, &v.Volume)
			if err != nil {
				return err
			}
		}
	}

	cntConfig := s.prepareContainerConfigFromInfo(info)
	hostConfig := s.prepareHostConfigFromInfo(info)

	err = s.pullDockerImage(ctx, cntConfig.Image)
	if err != nil {
		return err
	}

	res, err := s.docker.ContainerCreate(
		ctx,
		cntConfig,
		hostConfig,
		nil,
		nil,
		info.Container.DockerName,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.repo.UpdateContainerDockerID(ctx, id, res.ID)
	if err != nil {
		return err
	}

	return s.CreateTasks(ctx, &models.DockerTask{
		ContainerID: id,
		Action:      models.DTAStart,
	})
}

func (s *Service) createVolume(
	ctx context.Context,
	vol *models.Volume,
) error {
	res, err := s.docker.VolumeCreate(ctx, volume.CreateOptions{
		Name: s.generateVolumeName(vol),
		Labels: map[string]string{
			LabelPrefix + "installation_id": s.GetInstallationID(),
			LabelPrefix + "internal_id":     strconv.FormatInt(vol.ID, 10),
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	vol.DockerName = res.Name

	err = s.repo.UpdateVolume(ctx, vol)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) pullDockerImage(
	ctx context.Context,
	refStr string,
) error {
	body, err := s.docker.ImagePull(ctx, refStr, image.PullOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	defer body.Close()

	_, err = io.ReadAll(body)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) isContainerExist(
	ctx context.Context,
	id int64,
) (bool, error) {
	cnts, err := s.getDockerContainersRaw(ctx, structs.ContainersOptions{
		OnlyCurrentInstallation: true,
		InternalID:              id,
	})
	if err != nil {
		return false, err
	}

	return len(cnts) > 0, nil
}
