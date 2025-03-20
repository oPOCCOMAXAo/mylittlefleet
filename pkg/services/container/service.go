package container

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/settings"
	"github.com/pkg/errors"
)

type Service struct {
	docker   *client.Client
	settings *settings.Service

	// Unique ID for the installation. Required for the container to be able to identify itself.
	installationID string
}

type Config struct {
	InsideDocker bool `env:"INSIDE_DOCKER"` // Ignore. For internal use only.
}

func NewService(
	docker *client.Client,
	settings *settings.Service,
) *Service {
	return &Service{
		docker:   docker,
		settings: settings,
	}
}

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

func (s *Service) OnStart(ctx context.Context) error {
	err := s.initInstallationID(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetInstallationID() string {
	return s.installationID
}
