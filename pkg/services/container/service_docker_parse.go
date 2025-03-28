package container

import (
	"context"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/docker"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (s *Service) parseFullInfoFromSummary(
	ctx context.Context,
	summary *container.Summary,
) (*models.FullContainerInfo, error) {
	var (
		res models.FullContainerInfo
		err error
	)

	res.Container = s.parseContainerFromSummary(summary)
	res.Volumes = s.parseVolumesFromSummary(summary)

	inspect, err := s.docker.ContainerInspect(ctx, summary.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res.Ports = s.parsePortsFromInspect(&inspect)

	imageInspect, err := s.docker.ImageInspect(ctx, inspect.Image)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res.Envs = s.parseEnvsFromInspect(&inspect, &imageInspect)

	for _, vol := range res.Volumes {
		vol.ContainerID = res.Container.ID
	}

	for _, port := range res.Ports {
		port.ContainerID = res.Container.ID
	}

	for _, env := range res.Envs {
		env.ContainerID = res.Container.ID
	}

	return &res, nil
}

func (s *Service) parseContainerFromSummary(
	summary *container.Summary,
) *models.Container {
	var res models.Container

	res.DockerID = summary.ID

	res.DockerName = lo.CoalesceOrEmpty(summary.Names...)
	res.DockerName = strings.TrimPrefix(res.DockerName, "/")

	res.Paused = summary.State != docker.StateRunning
	res.Image, res.Tag = docker.ParseImage(summary.Image)
	res.Name = summary.Labels[LabelPrefix+"name"]

	{
		// Empty ID will be deleted.
		strID := summary.Labels[LabelPrefix+"internal_id"]
		if strID != "" {
			res.ID, _ = strconv.ParseInt(strID, 10, 64)
		}
	}

	return &res
}

func (s *Service) parseVolumesFromSummary(
	summary *container.Summary,
) []*models.VolumeDomain {
	res := make([]*models.VolumeDomain, 0, len(summary.Mounts))

	for _, mountPoint := range summary.Mounts {
		if mountPoint.Type != mount.TypeVolume {
			continue
		}

		res = append(res, &models.VolumeDomain{
			ContainerVolume: models.ContainerVolume{
				ContainerPath: mountPoint.Destination,
			},
			Volume: models.Volume{
				DockerName: mountPoint.Name,
			},
		})
	}

	return res
}

//nolint:cyclop
func (s *Service) parsePortsFromInspect(
	inspect *container.InspectResponse,
) []*models.ContainerPort {
	if inspect.HostConfig == nil {
		return nil
	}

	res := make([]*models.ContainerPort, 0, len(inspect.HostConfig.PortBindings))
	used := make(map[models.ContainePortUniqueKey]struct{})

	for port, bindings := range inspect.HostConfig.PortBindings {
		for _, binding := range bindings {
			if binding.HostIP == "" || binding.HostPort == "" {
				continue
			}

			temp := models.ContainerPort{
				ContainerPort: int64(port.Int()),
			}

			temp.HostPort, _ = strconv.ParseInt(binding.HostPort, 10, 64)
			if temp.HostPort == 0 {
				continue
			}

			switch binding.HostIP {
			case "127.0.0.1":
				temp.IsPublic = false
			case "0.0.0.0", "::":
				temp.IsPublic = true
			default:
				continue // unknown IP
			}

			uniqueKey := temp.UniqueKey()
			if _, ok := used[uniqueKey]; ok {
				continue
			}

			used[uniqueKey] = struct{}{}

			res = append(res, &temp)
		}
	}

	return res
}

func (s *Service) parseEnvsFromInspect(
	inspect *container.InspectResponse,
	imageInspect *image.InspectResponse,
) []*models.ContainerEnv {
	if inspect.Config == nil {
		return nil
	}

	res := make([]*models.ContainerEnv, 0, len(inspect.Config.Env))
	byName := make(map[string]*models.ContainerEnv, len(inspect.Config.Env))

	for _, env := range inspect.Config.Env {
		name, value := docker.ParseEnv(env)

		cEnv := &models.ContainerEnv{
			Name:         name,
			Value:        value,
			IsDefault:    false,
			DefaultValue: "",
		}

		res = append(res, cEnv)
		byName[name] = cEnv
	}

	if imageInspect.Config == nil {
		return res
	}

	for _, env := range imageInspect.Config.Env {
		name, value := docker.ParseEnv(env)

		cEnv, ok := byName[name]
		if !ok {
			temp := &models.ContainerEnv{
				Name:         name,
				Value:        value,
				IsDefault:    true,
				DefaultValue: value,
			}
			res = append(res, temp)
			byName[name] = temp

			continue
		}

		cEnv.IsDefault = true
		cEnv.DefaultValue = value
	}

	return res
}
