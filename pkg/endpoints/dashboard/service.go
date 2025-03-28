package dashboard

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/server"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/user"
)

type Service struct {
	user      *user.Service
	server    *server.Service
	container *container.Service
}

func NewService(
	user *user.Service,
	server *server.Service,
	container *container.Service,
) *Service {
	return &Service{
		user:      user,
		server:    server,
		container: container,
	}
}
