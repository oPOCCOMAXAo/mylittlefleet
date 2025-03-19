package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/views"
)

func (s *Service) SetupPage(ctx *gin.Context) {
	if s.auth.HasUsers() {
		ctx.Redirect(http.StatusFound, "/login")

		return
	}

	ctx.HTML(http.StatusOK, "", views.Setup())
}

type SetupRequest struct {
	Username string `form:"username,required"`
	Password string `form:"password,required"`
}

func (s *Service) Setup(ctx *gin.Context) {
	var req SetupRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "", views.Setup())

		return
	}

	user, err := s.auth.CreateUser(ctx.Request.Context(), auth.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "", views.Setup())

		return
	}

	err = s.auth.SetupSession(ctx, user)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "", views.Setup())

		return
	}

	ctx.Redirect(http.StatusFound, "/dashboard")
}
