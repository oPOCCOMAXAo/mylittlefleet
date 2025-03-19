package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/ginutils"
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
	auth *auth.Service,
	service *Service,
) {
	router.GET("/", auth.MiddlewareAuth, ginutils.StaticRedirect("/dashboard"))

	dashboard := router.Group("/dashboard")
	dashboard.Use(auth.MiddlewareAuth)

	dashboard.GET("", service.Dashboard)
}
