package event

import "github.com/opoccomaxao/mylittlefleet/pkg/services/event"

type Service struct {
	event *event.Service
}

func NewService(
	event *event.Service,
) *Service {
	return &Service{
		event: event,
	}
}
