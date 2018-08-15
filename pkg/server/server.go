package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
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
	e.Use(middleware.CORS())
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

			return h(c)
		}
	})

	// support batch requests
	e.Use(batchRequestsMiddleware)

	e.HTTPErrorHandler = errorHandler
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

func batchRequestsMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		myctx := c.Get("myctx")
		cc, ok := myctx.(*myCtx)
		if !ok {
			return errors.New("Could not find myctx")
		}

		// Request
		reqBody := []byte{}
		if c.Request().Body != nil { // Read
			reqBody, _ = ioutil.ReadAll(c.Request().Body)
		}
		isBatchRequests := func(msg json.RawMessage) bool {
			return msg[0] == '['
		}
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

		if !isBatchRequests(reqBody) {
			return h(c)
		}

		var rpcReqs []*eth.JSONRPCRequest
		if err := c.Bind(&rpcReqs); err != nil {
			return err
		}

		results := make([]*eth.JSONRPCResult, 0, len(rpcReqs))

		for _, req := range rpcReqs {
			result, err := callHttpHandler(cc, req)
			if err != nil {
				return err
			}

			results = append(results, result)
		}

		return c.JSON(http.StatusOK, results)
	}
}

func callHttpHandler(cc *myCtx, req *eth.JSONRPCRequest) (*eth.JSONRPCResult, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpreq := httptest.NewRequest(echo.POST, "/", ioutil.NopCloser(bytes.NewReader(reqBytes)))
	httpreq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newCtx := cc.Echo().NewContext(httpreq, rec)
	myCtx := &myCtx{
		Context:     newCtx,
		logger:      cc.logger,
		transformer: cc.transformer,
	}
	newCtx.Set("myctx", myCtx)

	if err = httpHandler(myCtx); err != nil {
		errorHandler(err, myCtx)
	}

	var result *eth.JSONRPCResult
	if err = json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		return nil, err
	}

	return result, nil
}
