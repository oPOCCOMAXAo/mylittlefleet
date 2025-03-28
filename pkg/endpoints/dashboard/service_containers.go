package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/views"
)

func (s *Service) ContainerList(ctx *gin.Context) {
	containers, err := s.container.GetContainersForDashboard(ctx.Request.Context())

	ctx.HTML(http.StatusOK, "", views.Dashboard(views.DashboardConfig{
		Page:       views.PageContainerList,
		Containers: containers,
		Error:      err,
	}))
}
