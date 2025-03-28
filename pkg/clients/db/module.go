package db

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/migrations"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("clients/db",
		fx.Provide(
			fx.Annotate(NewSQLite,
				fx.OnStart(migrations.Migrate),
			),
		),
	)
}
