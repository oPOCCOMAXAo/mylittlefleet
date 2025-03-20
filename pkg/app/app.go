package app

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/db"
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/docker"
	"github.com/opoccomaxao/mylittlefleet/pkg/config"
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints"
	"github.com/opoccomaxao/mylittlefleet/pkg/server"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/logger"
	serverSvc "github.com/opoccomaxao/mylittlefleet/pkg/services/server"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/settings"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/user"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	app := fx.New(
		cfg.Provide(),
		fx.Provide(NewCancelCause),
		fx.WithLogger(NewFxLogger),

		// modules.
		logger.Module(),
		db.Module(),
		server.Module(),
		auth.Module(),
		user.Module(),
		settings.Module(),
		serverSvc.Module(),
		docker.Module(),
		container.Module(),

		// invocations.
		endpoints.Invoke(),
	)
	app.Run()

	err = app.Err()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
