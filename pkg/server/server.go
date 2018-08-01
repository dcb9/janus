package server

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	address       string
	transformer   *transformer.Transformer
	qtumRPCClient *qtum.Qtum
	logger        log.Logger
	debug         bool
	echo          *echo.Echo
}

func New(
	qtumRPCClient *qtum.Qtum,
	transformer *transformer.Transformer,
	addr string,
	opts ...Option,
) (*Server, error) {
	p := &Server{
		logger:        log.NewNopLogger(),
		echo:          echo.New(),
		address:       addr,
		qtumRPCClient: qtumRPCClient,
		transformer:   transformer,
	}

	var err error
	for _, opt := range opts {
		if err = opt(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (s *Server) Start() error {
	e := s.echo
	e.HTTPErrorHandler = errorHandler
	e.Use(middleware.BodyDump(func(c echo.Context, req []byte, resp []byte) {
		myctx := c.Get("myctx")
		cc, ok := myctx.(*myCtx)
		if !ok {
			return
		}

		level.Debug(cc.logger).Log("msg", "body dump", "reqBody", req, "respBody", resp)
	}))

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &myCtx{
				Context:     c,
				logger:      s.logger,
				transformer: s.transformer,
			}

			c.Set("myctx", cc)

			var rpcReq *eth.JSONRPCRequest
			if err := c.Bind(&rpcReq); err != nil {
				return err
			}

			cc.rpcReq = rpcReq

			return h(c)
		}
	})
	e.HideBanner = true
	e.POST("/*", httpHandler)

	level.Warn(s.logger).Log("listen", s.address, "qtum_rpc", s.qtumRPCClient.URL, "msg", "proxy started")
	return e.Start(s.address)
}

type Option func(*Server) error

func SetLogger(l log.Logger) Option {
	return func(p *Server) error {
		p.logger = l
		return nil
	}
}

func SetDebug(debug bool) Option {
	return func(p *Server) error {
		p.debug = debug
		return nil
	}
}
