package container

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Invoke() fx.Option {
	return fx.Module("endpoints/dashboard",
		fx.Provide(NewService, fx.Private),
		fx.Invoke(RegisterEndpoints),
	)
}

func RegisterEndpoints(
	router gin.IRouter,
	service *Service,
) {
	router.GET("/installation_id", service.InstallationID)
}
