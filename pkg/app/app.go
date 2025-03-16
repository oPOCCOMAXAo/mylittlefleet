package app

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/db"
	"github.com/opoccomaxao/mylittlefleet/pkg/config"
	"github.com/opoccomaxao/mylittlefleet/pkg/server"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/logger"
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
	)
	app.Run()

	err = app.Err()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
