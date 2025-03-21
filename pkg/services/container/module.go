package container

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/repo"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/fxutils"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services/container",
		fxutils.ProvideWithHooks[*Service](NewService),
		fx.Provide(repo.NewRepo, fx.Private),
	)
}
