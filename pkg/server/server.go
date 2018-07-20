package server

import (
	"net/url"

	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/labstack/echo"
)

type Server struct {
	qtumRPC            *url.URL
	address            string
	transformerManager *transformer.Manager
	qtumClient         *qtum.Client
	logger             log.Logger
	debug              bool
}

func New(qtumRPC string, addr string, opts ...Option) (*Server, error) {
	opts = append(opts, setQtumRPC(qtumRPC), setAddress(addr))

	p := &Server{
		logger: log.NewNopLogger(),
	}

	var err error
	for _, opt := range opts {
		if err = opt(p); err != nil {
			return nil, err
		}
	}

	p.qtumClient, err = qtum.NewClient(
		qtumRPC,
		qtum.SetLogger(p.logger),
		qtum.SetDebug(p.debug),
	)
	if err != nil {
		return nil, err
	}

	p.transformerManager, err = transformer.NewManager(
		p.qtumClient,
		transformer.SetLogger(p.logger),
		transformer.SetDebug(p.debug),
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true
	e.POST("/*", s.httpHandler)
	e.HTTPErrorHandler = s.errorHandler

	level.Warn(s.logger).Log("listen", s.address, "qtum_rpc", s.qtumRPC, "msg", "proxy started")
	return e.Start(s.address)
}
