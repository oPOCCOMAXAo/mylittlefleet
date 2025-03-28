package hx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Invoke() fx.Option {
	return fx.Module("endpoints/hx",
		fx.Provide(NewService, fx.Private),
		fx.Invoke(RegisterEndpoints),
	)
}

func RegisterEndpoints(
	router gin.IRouter,
	service *Service,
) {
	route := router.Group("/hx")
	route.GET("/empty", service.Empty)
}
