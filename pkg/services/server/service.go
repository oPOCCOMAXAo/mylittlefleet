package server

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/settings"
)

type Service struct {
	settings *settings.Service

	nginxStatus models.ContainerStatus
}

func NewService(
	settings *settings.Service,
) *Service {
	return &Service{
		settings: settings,

		nginxStatus: models.CSError,
	}
}

func (s *Service) OnStart(ctx context.Context) error {
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

	res.NginxStatus = s.nginxStatus

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

	return nil
}
