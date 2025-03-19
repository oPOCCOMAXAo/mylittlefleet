package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/views"
)

func (s *Service) DashboardPage(ctx *gin.Context) {
	user, err := s.user.GetUserByID(
		ctx.Request.Context(),
		auth.CtxUserID.Get(ctx),
	)

	ctx.HTML(http.StatusOK, "", views.Dashboard(views.DashboardConfig{
		Page:  views.PageProfile,
		User:  user,
		Error: err,
	}))
}
