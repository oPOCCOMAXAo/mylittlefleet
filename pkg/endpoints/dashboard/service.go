package dashboard

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/services/server"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/user"
)

type Service struct {
	user   *user.Service
	server *server.Service
}

func NewService(
	user *user.Service,
	server *server.Service,
) *Service {
	return &Service{
		user:   user,
		server: server,
	}
}
