package server

import (
	"context"
	"log/slog"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/settings"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/contextutils"
)

type Service struct {
	logger    *slog.Logger
	settings  *settings.Service
	container *container.Service
}

func NewService(
	logger *slog.Logger,
	settings *settings.Service,
	container *container.Service,
) *Service {
	return &Service{
		logger:    logger,
		settings:  settings,
		container: container,
	}
}

func (s *Service) OnStart(ctx context.Context) error {
	err := s.ApplyConfig(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetServerConfig(
	ctx context.Context,
) (*models.ServerConfig, error) {
	var (
		res models.ServerConfig
		err error
	)

	res.ReverseProxyEnabled, err = s.settings.GetBool(ctx, "server:reverse_proxy_enabled")
	if err != nil {
		return nil, err
	}

	res.NginxStatus, err = s.container.GetContainerRuntimeStatusByName(ctx, container.ContainerNameReverseProxy)
	if err != nil {
		res.NginxStatus = models.CSError
	}

	return &res, nil
}

func (s *Service) SetServerConfig(
	ctx context.Context,
	cfg models.ServerConfig,
) error {
	err := s.settings.SetBool(ctx, "server:reverse_proxy_enabled", cfg.ReverseProxyEnabled)
	if err != nil {
		return err
	}

	s.ApplyConfigAsync(ctx)

	return nil
}

func (s *Service) ApplyConfig(ctx context.Context) error {
	cfg, err := s.GetServerConfig(ctx)
	if err != nil {
		return err
	}

	if cfg.ReverseProxyEnabled {
		err = s.container.StartContainerByName(ctx, container.ContainerNameReverseProxy)
		if err != nil {
			return err
		}
	} else {
		err = s.container.StopContainerByName(ctx, container.ContainerNameReverseProxy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ApplyConfigAsync(
	ctx context.Context,
) {
	ctx = contextutils.Detached(ctx)

	go func() {
		err := s.ApplyConfig(ctx)
		if err != nil {
			s.logger.ErrorContext(ctx, "ApplyConfig",
				slog.Any("error", err),
			)
		}
	}()
}
