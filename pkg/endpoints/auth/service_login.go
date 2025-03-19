package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/views"
)

func (s *Service) LoginPage(ctx *gin.Context) {
	if !s.auth.HasUsers() {
		ctx.Redirect(http.StatusFound, "/setup")

		return
	}

	ctx.HTML(http.StatusOK, "", views.Login())
}

type LoginRequest struct {
	Username string `form:"username,required"`
	Password string `form:"password,required"`
}

func (s *Service) Login(ctx *gin.Context) {
	var req LoginRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "", views.Login())

		return
	}

	user, err := s.auth.AuthUser(ctx.Request.Context(), auth.AuthUserParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		ctx.HTML(http.StatusUnauthorized, "", views.Login())

		return
	}

	err = s.auth.SetupSession(ctx, user)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "", views.Setup())

		return
	}

	ctx.Redirect(http.StatusFound, "/dashboard")
}
