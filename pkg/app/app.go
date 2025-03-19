package app

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/db"
	"github.com/opoccomaxao/mylittlefleet/pkg/config"
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints"
	"github.com/opoccomaxao/mylittlefleet/pkg/server"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/logger"
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
		logger.Module(),
		db.Module(),
		server.Module(),
		auth.Module(),
		user.Module(),
		endpoints.Invoke(),
	)
	app.Run()

	err = app.Err()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
