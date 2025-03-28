package hx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsHX(ctx *gin.Context) bool {
	return ctx.GetHeader("HX-Request") == "true"
}

func Redirect(ctx *gin.Context, url string) {
	if IsHX(ctx) {
		ctx.Header("HX-Redirect", url)
	} else {
		ctx.Redirect(http.StatusFound, url)
	}
}

func GetTarget(ctx *gin.Context) string {
	return ctx.GetHeader("HX-Target")
}
