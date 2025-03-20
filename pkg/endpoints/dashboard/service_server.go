package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/hx"
	"github.com/opoccomaxao/mylittlefleet/pkg/views"
)

func (s *Service) ServerPage(ctx *gin.Context) {
	cfg, err := s.server.GetServerConfig(ctx.Request.Context())

	ctx.HTML(http.StatusOK, "", views.Dashboard(views.DashboardConfig{
		Page:   views.PageServer,
		Server: cfg,
		Error:  err,
	}))
}

func (s *Service) ServerEditPage(ctx *gin.Context) {
	cfg, err := s.server.GetServerConfig(ctx.Request.Context())

	ctx.HTML(http.StatusOK, "", views.Dashboard(views.DashboardConfig{
		Page:   views.PageServerEdit,
		Server: cfg,
		Error:  err,
	}))
}

type ServerUpdateForm struct {
	ReverseProxyEnabled bool `form:"reverse_proxy_enabled"`
}

func (s *Service) ServerUpdate(ctx *gin.Context) {
	var req ServerUpdateForm

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "", views.Dashboard(views.DashboardConfig{
			Page:  views.PageServerEdit,
			Error: err,
		}))

		return
	}

	err = s.server.SetServerConfig(ctx.Request.Context(), models.ServerConfig{
		ReverseProxyEnabled: req.ReverseProxyEnabled,
	})
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "", views.Dashboard(views.DashboardConfig{
			Page:  views.PageServerEdit,
			Error: err,
		}))

		return
	}

	hx.Redirect(ctx, "/dashboard/server")
}
