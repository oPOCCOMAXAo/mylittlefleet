package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/hx"
)

func (s *Service) Logout(ctx *gin.Context) {
	err := s.auth.ClearSession(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	hx.Redirect(ctx, "/")
}
