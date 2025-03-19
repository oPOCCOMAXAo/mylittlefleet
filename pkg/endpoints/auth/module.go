package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Invoke() fx.Option {
	return fx.Module("endpoints/auth",
		fx.Provide(NewService, fx.Private),
		fx.Invoke(RegisterEndpoints),
	)
}

func RegisterEndpoints(
	router gin.IRouter,
	service *Service,
) {
	router.GET("/setup", service.SetupPage)
	router.POST("/setup", service.Setup)
	router.GET("/login", service.LoginPage)
	router.POST("/login", service.Login)
	router.POST("/logout", service.Logout)
}
