package server

import (
	stdLog "log"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/go-kit/kit/log/level"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func httpHandler(c echo.Context) error {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if !ok {
		return errors.New("Could not find myctx")
	}

	var rpcReq *eth.JSONRPCRequest
	if err := c.Bind(&rpcReq); err != nil {
		return err
	}

	cc.rpcReq = rpcReq

	level.Debug(cc.logger).Log("msg", "before call transformer#Transform")
	result, err := cc.transformer.Transform(rpcReq)
	level.Debug(cc.logger).Log("msg", "after call transformer#Transform")

	if err != nil {
		err1 := errors.Cause(err)
		if err != err1 {
			level.Error(cc.logger).Log("err", err.Error())
			return cc.JSONRPCError(&eth.JSONRPCError{
				Code:    100,
				Message: err1.Error(),
			})
		}

		return err
	}

	return cc.JSONRPCResult(result)
}

func errorHandler(err error, c echo.Context) {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if ok {
		level.Error(cc.logger).Log("err", err.Error())
		if err := cc.JSONRPCError(&eth.JSONRPCError{
			Code:    100,
			Message: err.Error(),
		}); err != nil {
			level.Error(cc.logger).Log("msg", "reply to client", "err", err.Error())
		}
		return
	}

	stdLog.Println("errorHandler", err.Error())
}
