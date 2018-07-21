package server

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/go-kit/kit/log"
)

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

func setAddress(addr string) Option {
	return func(p *Server) error {
		p.address = addr
		return nil
	}
}

func setQtumRPC(r string) Option {
	return func(p *Server) error {
		if r == "" {
			return errors.New("Please set QTUM_RPC to qtumd's RPC URL")
		}

		qtumRPC, err := url.Parse(r)
		if err != nil {
			return errors.New(fmt.Sprintf("QTUM_RPC URL: %s", r))
		}

		if qtumRPC.User == nil {
			return errors.New(fmt.Sprintf("QTUM_RPC URL (must specify user & password): %s", r))
		}

		p.qtumRPC = qtumRPC

		return nil
	}
}
