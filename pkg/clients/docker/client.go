package docker

import (
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type Config struct {
	Host string `env:"HOST" envDefault:"unix:///var/run/docker.sock"` // Docker daemon socket. Must be local.
}

func NewClient(config Config) (*client.Client, error) {
	opts := []client.Opt{
		client.WithAPIVersionNegotiation(),
	}
	if config.Host != "" {
		opts = append(opts, client.WithHost(config.Host))
	}

	client, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client, nil
}
