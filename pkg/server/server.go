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
	echo               *echo.Echo
}

// FIXME: define a StartServer/main function to do dependency injection. Constructor should be simple.

func New(qtumRPC string, addr string, opts ...Option) (*Server, error) {
	opts = append(opts, setQtumRPC(qtumRPC), setAddress(addr))

	p := &Server{
		logger: log.NewNopLogger(),
		echo:   echo.New(),
	}

	var err error
	for _, opt := range opts {
		if err = opt(p); err != nil {
			return nil, err
		}
	}

	// FIXME: dependency injection: pass client an transformer into constructor as parameters

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
	e := s.echo
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &myCtx{
				Context: c,
				server:  s,
				logger:  s.logger,
			}
			c.Set("myctx", cc)
			return h(cc)
		}
	})

	e.HideBanner = true
	e.POST("/*", s.myCtxHandler(httpHandler))
	e.HTTPErrorHandler = errorHandler

	level.Warn(s.logger).Log("listen", s.address, "qtum_rpc", s.qtumRPC, "msg", "proxy started")
	return e.Start(s.address)
}
