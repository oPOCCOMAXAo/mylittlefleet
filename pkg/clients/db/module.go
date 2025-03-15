package db

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/migrations"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func Module() fx.Option {
	return fx.Module("clients/db",
		fx.Provide(newModule),
	)
}

type moduleParams struct {
	fx.In
	fx.Lifecycle

	Config Config
}

type moduleResults struct {
	fx.Out

	DB *gorm.DB
}

func newModule(
	params moduleParams,
) (moduleResults, error) {
	var (
		res moduleResults
		err error
	)

	res.DB, err = NewSQLite(params.Config)
	if err != nil {
		return res, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return migrations.Migrate(ctx, res.DB)
		},
	})

	return res, nil
}
