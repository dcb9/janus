package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type proxy struct {
	qtumRPC *url.URL
	address string
	logger  log.Logger
}

func New(qtumRPC string, addr string, opts ...Option) (*proxy, error) {
	opts = append(opts, setQtumRPC(qtumRPC), setAddress(addr))

	p := &proxy{
		logger: log.NewNopLogger(),
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *proxy) Start() error {
	rp := httputil.NewSingleHostReverseProxy(p.qtumRPC)
	rp.Transport = &Transport{
		reqTransformers:  transformer.RequestTransformers,
		respTransformers: transformer.ResponseTransformers,
		logger:           p.logger,
		userInfo:         p.qtumRPC.User,
	}
	level.Warn(p.logger).Log("listen", p.address, "qtum_rpc", p.qtumRPC, "msg", "proxy started")
	return http.ListenAndServe(p.address, rp)
}
