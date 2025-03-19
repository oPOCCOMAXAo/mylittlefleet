package auth

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth/repo"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/fxutils"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services/auth",
		fxutils.ProvideWithHooks[*Service](NewService),
		fx.Provide(repo.NewRepo, fx.Private),
	)
}
