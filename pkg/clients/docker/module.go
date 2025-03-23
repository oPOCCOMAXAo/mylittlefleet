package docker

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("clients/docker",
		fx.Provide(NewClient),
	)
}
