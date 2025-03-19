package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) Logout(ctx *gin.Context) {
	err := s.auth.ClearSession(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Status(http.StatusOK)
}
