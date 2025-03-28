package event

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/fxutils"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services/event",
		fxutils.ProvideWithHooks[*Service](New),
	)
}
