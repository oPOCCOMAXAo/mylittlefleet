package server

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/fxutils"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("server",
		fxutils.ProvideWithHooks[*Server](New),
		fx.Provide((*Server).Router),
		fx.Provide((*Server).Engine),
	)
}
