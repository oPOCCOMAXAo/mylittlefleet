package auth

import "github.com/opoccomaxao/mylittlefleet/pkg/services/auth"

type Service struct {
	auth *auth.Service
}

func NewService(
	auth *auth.Service,
) *Service {
	return &Service{
		auth: auth,
	}
}
