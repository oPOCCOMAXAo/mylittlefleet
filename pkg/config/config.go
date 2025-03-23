package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/db"
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/docker"
	"github.com/opoccomaxao/mylittlefleet/pkg/server"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/logger"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// App configuration, environment variables.
//
//go:generate envdoc -output ./../../.example.env -dir ../ -target caarlos0 -types Config -files ./config/config.go -format dotenv
type Config struct {
	fx.Out

	Extra     Extra            `envPrefix:""`
	Logger    logger.Config    `envPrefix:"LOGGER_"`
	Server    server.Config    `envPrefix:"SERVER_"`
	DB        db.Config        `envPrefix:"DB_"`
	Auth      auth.Config      `envPrefix:"AUTH_"`
	Docker    docker.Config    `envPrefix:"DOCKER_"`
	Container container.Config `envPrefix:"CONTAINER_"`
}

type Extra struct {
	StartTimeout time.Duration `env:"START_TIMEOUT"` // for debugging purposes.
}

func New() (Config, error) {
	var res Config

	err := env.ParseWithOptions(&res, env.Options{
		UseFieldNameByDefault: false,
		RequiredIfNoDef:       false,
	})
	if err != nil {
		return res, errors.WithStack(err)
	}

	return res, nil
}

func (c Config) Provide() fx.Option {
	opts := []fx.Option{
		fx.Provide(func() Config { return c }),
	}

	if c.Extra.StartTimeout > 0 {
		opts = append(opts, fx.StartTimeout(c.Extra.StartTimeout))
	}

	return fx.Options(opts...)
}
