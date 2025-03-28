package container

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/repo"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/event"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/settings"
)

type Service struct {
	mu                 sync.Mutex
	runCtx             context.Context //nolint:containedctx
	runCancel          context.CancelFunc
	chanSyncWithDocker chan struct{}
	chanTasks          chan struct{}

	repo     *repo.Repo
	logger   *slog.Logger
	httpCli  *http.Client
	docker   *client.Client
	settings *settings.Service
	event    *event.Service

	// Unique ID for the installation. Required for the container to be able to identify itself.
	installationID string
}

type Config struct {
	InsideDocker bool `env:"INSIDE_DOCKER"` // Ignore. For internal use only.
}

//nolint:mnd
func NewService(
	repo *repo.Repo,
	logger *slog.Logger,
	docker *client.Client,
	settings *settings.Service,
	event *event.Service,
) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "container")),
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
		docker:   docker,
		settings: settings,
		event:    event,

		runCtx:             context.Background(),
		chanSyncWithDocker: make(chan struct{}, 1),
		chanTasks:          make(chan struct{}, 1),
	}
}

func (s *Service) OnStart(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.runCancel != nil {
		return nil
	}

	s.runCtx, s.runCancel = context.WithCancel(context.Background())

	var err error

	for _, init := range []func(context.Context) error{
		s.initInstallationID,
		s.initInternalContainers,
	} {
		err = init(ctx)
		if err != nil {
			return err
		}
	}

	go s.serveSyncWithDocker()
	go s.serveTasks()

	return nil
}

func (s *Service) OnStop(_ context.Context) error {
	if s.runCancel != nil {
		s.runCancel()
	}

	return nil
}
