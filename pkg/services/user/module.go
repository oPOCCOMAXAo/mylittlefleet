package user

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/services/user/repo"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services/user",
		fx.Provide(NewService),
		fx.Provide(repo.NewRepo, fx.Private),
	)
}
