package settings

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/services/settings/repo"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services/settings",
		fx.Provide(NewService),
		fx.Provide(repo.NewRepo, fx.Private),
	)
}
