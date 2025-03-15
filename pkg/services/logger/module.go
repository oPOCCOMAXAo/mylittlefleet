package logger

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services/logger",
		fx.Provide(New),
		fx.Decorate(DecorateGormDB),
	)
}
