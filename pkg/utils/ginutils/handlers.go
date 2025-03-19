package ginutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StaticRedirect(
	path string,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, path)
	}
}
