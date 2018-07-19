package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type proxy struct {
	qtumRPC            *url.URL
	address            string
	transformerManager *transformer.Manager
	qtumClient         *qtum.Client
	logger             log.Logger
	debug              bool
}

func New(qtumRPC string, addr string, opts ...Option) (*proxy, error) {
	opts = append(opts, setQtumRPC(qtumRPC), setAddress(addr))

	p := &proxy{
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

func (p *proxy) Start() error {
	rp := httputil.NewSingleHostReverseProxy(p.qtumRPC)
	rp.Transport = p

	level.Warn(p.logger).Log("listen", p.address, "qtum_rpc", p.qtumRPC, "msg", "proxy started")
	return http.ListenAndServe(p.address, rp)
}