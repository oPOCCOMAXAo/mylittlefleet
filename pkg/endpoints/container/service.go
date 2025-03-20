package container

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container"
)

type Service struct {
	container *container.Service
}

func NewService(
	container *container.Service,
) *Service {
	return &Service{
		container: container,
	}
}
