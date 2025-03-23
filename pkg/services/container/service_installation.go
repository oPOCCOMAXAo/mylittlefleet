package container

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/docker"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/structs"
	"github.com/pkg/errors"
)

func (s *Service) initInstallationID(ctx context.Context) error {
	var err error

	s.installationID, err = s.settings.Get(ctx, "container:installation_id")
	if err != nil {
		return err
	}

	if s.installationID != "" {
		return nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return errors.WithStack(err)
	}

	s.installationID = id.String()

	err = s.settings.Set(ctx, "container:installation_id", s.installationID)
	if err != nil {
		return err
	}

	return nil
}

//nolint:mnd
func (s *Service) initInternalContainers(ctx context.Context) error {
	containers := []*structs.FullContainerInfo{
		{
			Container: &models.Container{
				Name:     ContainerNameSelf,
				Internal: true,
			},
		},
		{
			Container: &models.Container{
				Name:     ContainerNameReverseProxy,
				Image:    "nginx",
				Tag:      "stable-alpine",
				Internal: true,
			},
			Ports: []*models.ContainerPort{
				{
					ContainerPort: 80,
					HostPort:      80,
					IsPublic:      true,
				},
				{
					ContainerPort: 443,
					HostPort:      443,
					IsPublic:      true,
				},
			},
		},
	}

	for _, container := range containers {
		err := s.EnsureFullContainerSettings(ctx, container)
		if err != nil {
			return errors.Wrapf(err, "container: %s", container.Container.Name)
		}
	}

	return nil
}

func (s *Service) GetInstallationID() string {
	return s.installationID
}

func (s *Service) FindCurrentInstallation(
	ctx context.Context,
) (*container.Summary, error) {
	containers, err := s.docker.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, container := range containers {
		if container.State != docker.StateRunning {
			continue
		}

		isCurrent, err := s.checkIfContainerIsCurrentInstallation(ctx, &container)
		if err != nil {
			return nil, err
		}

		if isCurrent {
			return &container, nil
		}
	}

	return nil, errors.WithStack(models.ErrNotFound)
}

func (s *Service) checkIfContainerIsCurrentInstallation(
	ctx context.Context,
	container *container.Summary,
) (bool, error) {
	ips := make([]string, 0, len(container.NetworkSettings.Networks))

	for _, network := range container.NetworkSettings.Networks {
		if network.IPAddress == "" {
			continue
		}

		ips = append(ips, network.IPAddress)
	}

	if len(ips) == 0 {
		return false, nil
	}

	for _, port := range container.Ports {
		if port.PrivatePort == 0 {
			continue
		}

		for _, ip := range ips {
			otherID, err := s.getOtherInstallationID(ctx, ip+":"+strconv.FormatUint(uint64(port.PrivatePort), 10))
			if err != nil && !errors.Is(err, models.ErrNotFound) {
				return false, err
			}

			if otherID == s.installationID {
				return true, nil
			}
		}
	}

	return false, nil
}

func (s *Service) getOtherInstallationID(
	ctx context.Context,
	host string,
) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://"+host+"/installation_id",
		http.NoBody,
	)
	if err != nil {
		return "", errors.WithStack(err)
	}

	res, err := s.httpCli.Do(req)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return "", errors.WithStack(models.ErrNotFound)
		}

		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			if strings.Contains(urlErr.Error(), "malformed HTTP response") {
				return "", errors.WithStack(models.ErrNotFound)
			}
		}

		return "", errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.WithStack(models.ErrNotFound)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(body), nil
}
