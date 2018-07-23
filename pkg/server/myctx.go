package server

import (
	"encoding/json"
	"net/http"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type myCtx struct {
	echo.Context
	rpcReq *rpc.JSONRPCRequest
	logger log.Logger
	server *Server
}

func (c *myCtx) JSONRPCResult(result interface{}) error {
	bytes, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "myCtx#JSONRPCResult")
	}
	return c.Context.JSON(http.StatusOK, eth.NewJSONRPCResult(c.rpcReq.ID, bytes, nil))
}

func (c *myCtx) JSONRPCError(err *rpc.JSONRPCError) error {
	return c.Context.JSON(http.StatusInternalServerError, eth.NewJSONRPCResult(c.rpcReq.ID, nil, err))
}

func (s *Server) myCtxHandler(h func(*myCtx) (result interface{}, err error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var rpcReq rpc.JSONRPCRequest
		if err := c.Bind(&rpcReq); err != nil {
			return err
		}

		c = &myCtx{
			Context: c,
			rpcReq:  &rpcReq,
			logger:  s.logger,
			server:  s,
		}

		cc, _ := c.(*myCtx)
		result, err := h(cc)
		if err != nil {
			return err
		}
		return cc.JSONRPCResult(result)
	}
}
