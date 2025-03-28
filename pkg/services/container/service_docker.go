package container

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/docker"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/structs"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (s *Service) serveSyncWithDocker() {
	s.RunSyncWithDocker()

	for {
		select {
		case <-s.runCtx.Done():
			return
		case <-s.chanSyncWithDocker:
			err := s.syncWithDocker(s.runCtx)
			if err != nil {
				s.logger.Error("syncWithDocker", slog.Any("error", err))
			}
		}
	}
}

func (s *Service) RunSyncWithDocker() {
	select {
	case s.chanSyncWithDocker <- struct{}{}:
	default:
	}
}

func (s *Service) syncWithDocker(ctx context.Context) error {
	current, err := s.findCurrentInstallation(ctx)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	err = s.syncCurrentInstallation(ctx, current)
	if err != nil {
		return err
	}

	err = s.syncContainers(ctx)
	if err != nil {
		return err
	}

	return nil
}

// dockerContainer may be nil when container is not running.
func (s *Service) syncCurrentInstallation(
	ctx context.Context,
	dockerContainer *container.Summary,
) error {
	oldContainer, err := s.GetContainerByName(ctx, ContainerNameSelf)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	if oldContainer == nil {
		return errors.WithStack(models.ErrFlowBroken)
	}

	if dockerContainer == nil {
		oldContainer.Paused = true
	} else {
		oldContainer.DockerID = dockerContainer.ID
		oldContainer.DockerName = lo.CoalesceOrEmpty(dockerContainer.Names...)
		oldContainer.Image, oldContainer.Tag = docker.ParseImage(dockerContainer.Image)
	}

	return s.repo.UpdateContainer(ctx, oldContainer)
}

//nolint:funlen,prealloc
func (s *Service) syncContainers(ctx context.Context) error {
	var (
		params structs.DiffContainersListParams
		err    error
	)

	params.Runtime, err = s.getDockerContainers(ctx, structs.ContainersOptions{
		OnlyCurrentInstallation: true,
	})
	if err != nil {
		return err
	}

	params.Storage, err = s.getAllFullContainerInfos(ctx)
	if err != nil {
		return err
	}

	diff := s.diffContainers(params)

	for _, container := range diff.StorageUpdate {
		err = s.SaveFullContainerSettings(ctx, container)
		if err != nil {
			return err
		}
	}

	var tasks []*models.DockerTask

	for _, container := range diff.DockerStop {
		tasks = append(tasks, &models.DockerTask{
			ContainerID: container.Container.ID,
			Action:      models.DTAStop,
		})
	}

	for _, container := range diff.DockerStart {
		tasks = append(tasks, &models.DockerTask{
			ContainerID: container.Container.ID,
			Action:      models.DTAStart,
		})
	}

	for _, container := range diff.DockerCreate {
		tasks = append(tasks, &models.DockerTask{
			ContainerID: container.Container.ID,
			Action:      models.DTACreate,
		})
	}

	for _, container := range diff.DockerUpdate {
		tasks = append(tasks,
			&models.DockerTask{
				ContainerID: container.Container.ID,
				Action:      models.DTADelete,
			},
			&models.DockerTask{
				ContainerID: container.Container.ID,
				Action:      models.DTACreate,
			},
		)
	}

	err = s.CreateTasks(ctx, tasks...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) getDockerContainersRaw(
	ctx context.Context,
	options structs.ContainersOptions,
) ([]container.Summary, error) {
	opts := container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(),
	}

	if options.OnlyCurrentInstallation {
		opts.Filters.Add("label", LabelPrefix+"installation_id="+s.GetInstallationID())
	}

	if options.InternalID != 0 {
		opts.Filters.Add("label", LabelPrefix+"internal_id="+strconv.FormatInt(options.InternalID, 10))
	}

	if options.DockerID != "" {
		opts.Filters.Add("id", options.DockerID)
	}

	if options.DockerName != "" {
		opts.Filters.Add("name", options.DockerName)
	}

	dockerInfo, err := s.docker.ContainerList(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return dockerInfo, nil
}

func (s *Service) getDockerContainers(
	ctx context.Context,
	options structs.ContainersOptions,
) ([]*structs.FullContainerInfo, error) {
	dockerInfo, err := s.getDockerContainersRaw(ctx, options)
	if err != nil {
		return nil, err
	}

	containers := make([]*structs.FullContainerInfo, 0, len(dockerInfo))

	for _, info := range dockerInfo {
		container, err := s.parseFullInfoFromSummary(ctx, &info)
		if err != nil {
			return nil, err
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func (s *Service) prepareContainerConfigFromInfo(
	info *structs.FullContainerInfo,
) *container.Config {
	var res container.Config

	res.Image = info.Container.Image + ":" + info.Container.Tag

	res.Labels = map[string]string{
		LabelPrefix + "installation_id": s.GetInstallationID(),
		LabelPrefix + "internal_id":     strconv.FormatInt(info.Container.ID, 10),
		"com.docker.compose.project":    ContainerNamePrefix + s.GetInstallationID(),
	}

	res.Env = make([]string, 0, len(info.Envs))

	for _, env := range info.Envs {
		if env.Value == env.DefaultValue {
			continue
		}

		res.Env = append(res.Env, env.Name+"="+env.Value)
	}

	res.Volumes = make(map[string]struct{}, len(info.Volumes))

	for _, v := range info.Volumes {
		if v.Volume.DockerName != "" {
			res.Volumes[v.Volume.DockerName] = struct{}{}
		}
	}

	res.ExposedPorts = make(nat.PortSet, len(info.Ports))

	for _, port := range info.Ports {
		natPort := nat.Port(strconv.FormatInt(port.ContainerPort, 10) + "/tcp")
		res.ExposedPorts[natPort] = struct{}{}
	}

	return &res
}

func (s *Service) prepareHostConfigFromInfo(
	info *structs.FullContainerInfo,
) *container.HostConfig {
	var res container.HostConfig

	res.Binds = make([]string, 0, len(info.Volumes))

	for _, v := range info.Volumes {
		if v.Volume.DockerName != "" {
			res.Binds = append(res.Binds, v.Volume.DockerName+":"+v.ContainerPath)
		}
	}

	res.PortBindings = make(map[nat.Port][]nat.PortBinding, len(info.Ports))

	for _, port := range info.Ports {
		natPort := nat.Port(strconv.FormatInt(port.ContainerPort, 10) + "/tcp")

		binding := nat.PortBinding{
			HostPort: strconv.FormatInt(port.HostPort, 10),
		}

		if port.IsPublic {
			binding.HostIP = "0.0.0.0"
		} else {
			binding.HostIP = "127.0.0.1"
		}

		res.PortBindings[natPort] = append(res.PortBindings[natPort], binding)
	}

	res.RestartPolicy = container.RestartPolicy{
		Name: container.RestartPolicyUnlessStopped,
	}

	return &res
}

func (s *Service) GetContainerRuntimeStatusByName(
	ctx context.Context,
	name string,
) (models.ContainerStatus, error) {
	containerID, err := s.repo.GetContainerIDByName(ctx, name)
	if err != nil {
		return models.CSError, err
	}

	res, err := s.getDockerContainersRaw(ctx, structs.ContainersOptions{
		OnlyCurrentInstallation: true,
		InternalID:              containerID,
	})
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return models.CSStopped, nil
	}

	return models.CSStatic.FromDockerState(res[0].State), nil
}
