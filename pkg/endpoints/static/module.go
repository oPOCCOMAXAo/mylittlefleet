package static

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/assets"
	"go.uber.org/fx"
)

func Invoke() fx.Option {
	return fx.Module("endpoints/static",
		fx.Invoke(RegisterEndpoints),
	)
}

func RegisterEndpoints(
	router gin.IRouter,
) {
	router.StaticFS("/assets", http.FS(assets.Assets))
}
