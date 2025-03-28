package hx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) Empty(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
