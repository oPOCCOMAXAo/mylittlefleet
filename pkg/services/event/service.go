package event

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
)

type Service struct {
	mu        sync.Mutex
	runCtx    context.Context //nolint:containedctx
	runCancel context.CancelFunc

	containerListener *Listener[models.ContainerEvent]
}

func New() *Service {
	return &Service{
		containerListener: NewListener[models.ContainerEvent](),
	}
}

func (s *Service) OnStart(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.runCtx != nil {
		return errors.WithStack(models.ErrFlowBroken)
	}

	s.runCtx, s.runCancel = context.WithCancel(ctx)

	go s.containerListener.Serve()

	return nil
}

func (s *Service) OnStop(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.runCancel != nil {
		s.runCancel()
	}

	s.containerListener.Close()

	return nil
}

func (s *Service) MiddlewarePrepareSSE(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")
}

// WatchContainer returns a channel that will receive container events.
// The channel will be closed when the context is done.
func (s *Service) WatchContainer(ctx *gin.Context) <-chan models.ContainerEvent {
	client := s.containerListener.NewClient()

	go func() {
		<-ctx.Request.Context().Done()

		s.containerListener.CloseClient(client)
	}()

	return client
}

func (s *Service) NotifyContainerEvent(
	_ context.Context,
	container models.ContainerEvent,
) error {
	s.containerListener.Notify(container)

	return nil
}
