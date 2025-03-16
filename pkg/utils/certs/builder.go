package certs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
)

// Builder is a helper to build certificates
//
// Don't use this directly, use certs.NewBuilder instead.
type Builder struct {
	template x509.Certificate
}

type Result struct {
	Certificate tls.Certificate
	CertPEM     []byte
	KeyPEM      []byte
}

func NewBuilder() *Builder {
	return &Builder{
		template: x509.Certificate{
			Issuer: pkix.Name{
				CommonName:         "localhost",
				Organization:       []string{"mylittlefleet"},
				OrganizationalUnit: []string{"mylittlefleet"},
			},
			Subject: pkix.Name{
				CommonName:         "localhost",
				Organization:       []string{"mylittlefleet"},
				OrganizationalUnit: []string{"mylittlefleet"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(1, 0, 0),
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature,
			BasicConstraintsValid: true,
			Version:               tls.VersionTLS13,
		},
	}
}

// AddHost adds a host to the certificate.
//
// It can be an IP address or a DNS name.
//
// If the host is empty, it will be ignored.
func (b *Builder) AddHost(host string) *Builder {
	if host == "" {
		return b
	}

	ip := net.ParseIP(host)
	if ip != nil {
		b.template.IPAddresses = append(b.template.IPAddresses, ip)
	} else {
		b.template.DNSNames = append(b.template.DNSNames, host)
	}

	return b
}

//nolint:mnd
func (b *Builder) Build() (*Result, error) {
	var err error

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	b.template.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&b.template,
		&b.template,
		&priv.PublicKey,
		priv,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var res Result

	res.CertPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res.KeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	res.Certificate, err = tls.X509KeyPair(res.CertPEM, res.KeyPEM)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &res, nil
}
