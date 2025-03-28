package endpoints

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints/container"
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints/dashboard"
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints/event"
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints/hx"
	"github.com/opoccomaxao/mylittlefleet/pkg/endpoints/static"
	"go.uber.org/fx"
)

func Invoke() fx.Option {
	return fx.Module("endpoints",
		static.Invoke(),
		hx.Invoke(),
		auth.Invoke(),
		dashboard.Invoke(),
		container.Invoke(),
		event.Invoke(),
	)
}
