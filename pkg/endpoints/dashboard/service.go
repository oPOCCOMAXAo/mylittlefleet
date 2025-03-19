package dashboard

import "github.com/opoccomaxao/mylittlefleet/pkg/services/user"

type Service struct {
	user *user.Service
}

func NewService(
	user *user.Service,
) *Service {
	return &Service{
		user: user,
	}
}
