package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/ginutils"
	pkgerr "github.com/pkg/errors"
)

type Server struct {
	config Config
	engine *gin.Engine
	http   *http.Server
	https  *http.Server
}

type Config struct {
	PortHTTP  string `env:"PORT_HTTP"  envDefault:"8080"` // HTTP port.
	PortHTTPS string `env:"PORT_HTTPS" envDefault:"8443"` // HTTPS port.
}

func New(
	config Config,
) *Server {
	res := &Server{
		config: config,
	}

	res.initEngine()
	res.initServer()

	return res
}

func (s *Server) initEngine() {
	s.engine = gin.New()
	_ = s.engine.SetTrustedProxies(nil)
	s.engine.Use(gin.Recovery())
	s.engine.HTMLRender = ginutils.NewTemplRender(s.engine.HTMLRender)
}

//nolint:mnd
func (s *Server) initServer() {
	s.http = &http.Server{
		Addr:              ":" + s.config.PortHTTP,
		Handler:           s.engine,
		ReadHeaderTimeout: 15 * time.Second,
	}
	s.https = &http.Server{
		Addr:              ":" + s.config.PortHTTPS,
		Handler:           s.engine,
		ReadHeaderTimeout: 15 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
}

func (s *Server) OnStart(
	ctx context.Context,
	cancelCauseFunc context.CancelCauseFunc,
) error {
	err := s.initCert(ctx)
	if err != nil {
		return err
	}

	go func() {
		err := s.http.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancelCauseFunc(err)
		}
	}()

	go func() {
		err := s.https.ListenAndServeTLS("", "")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancelCauseFunc(err)
		}
	}()

	return nil
}

func (s *Server) OnStop(ctx context.Context) error {
	errHTTP := s.http.Shutdown(ctx)
	errHTTPS := s.https.Shutdown(ctx)

	return pkgerr.WithStack(errors.Join(errHTTP, errHTTPS))
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Router() gin.IRouter {
	return s.engine
}
