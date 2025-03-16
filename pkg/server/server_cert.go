package server

import (
	"context"
	"crypto/tls"

	"github.com/opoccomaxao/mylittlefleet/pkg/utils/certs"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/netutils"
)

func (s *Server) initCert(ctx context.Context) error {
	cert, err := s.createCert(ctx)
	if err == nil {
		s.https.TLSConfig.Certificates = []tls.Certificate{cert}

		return nil
	}

	return err
}

func (*Server) createCert(ctx context.Context) (tls.Certificate, error) {
	publicIP, _ := netutils.GetCurrentPublicIP(ctx)

	certs, err := certs.NewBuilder().
		AddHost("127.0.0.1").
		AddHost("localhost").
		AddHost("*.localhost").
		AddHost(netutils.IPHostOrEmpty(publicIP)).
		Build()
	if err != nil {
		return tls.Certificate{}, err
	}

	return certs.Certificate, nil
}
