package event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/ginutils"
	"github.com/opoccomaxao/mylittlefleet/pkg/views"
)

func (s *Service) EventContainer(ctx *gin.Context) {
	ch := s.event.WatchContainer(ctx)

	for message := range ch {
		ginutils.RenderSSE(ctx, "container", views.Toast(views.ToastParams{
			Title: "Container update",
			Message: fmt.Sprintf("ID:%d\nName: %s\nStatus: %s",
				message.ID,
				message.Name,
				message.Status.String(),
			),
			Time:   message.Time,
			IsCode: true,
		}))
		ginutils.RenderSSE(ctx, "container-status", views.ContainerStatusBadge(views.ContainerStatusBadgeConfig{
			Status:        message.Status,
			ContainerName: message.Name,
			IsSSE:         true,
		}))
	}
}
