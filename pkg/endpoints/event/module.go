package event

import (
	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/event"
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
	event *event.Service,
) {
	route := router.Group("/event")
	route.Use(auth.MiddlewareAuth, event.MiddlewarePrepareSSE)

	route.GET("/container", service.EventContainer)
}
