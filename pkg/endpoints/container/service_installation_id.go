package container

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) InstallationID(ctx *gin.Context) {
	ctx.String(http.StatusOK, s.container.GetInstallationID())
}
