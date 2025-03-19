package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) Dashboard(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
