package server

import (
	"log"

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

	level.Debug(cc.logger).Log("msg", "before call transformer#Transform")
	result, err := cc.transformer.Transform(cc.rpcReq)
	level.Debug(cc.logger).Log("msg", "after call transformer#Transform")
	if err != nil {
		return err
	}

	return cc.JSONRPCResult(result)
}

func errorHandler(err error, c echo.Context) {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if ok {
		err1 := errors.Cause(err)
		if err != err1 {
			level.Error(cc.logger).Log("err", err.Error())
		}
		cc.JSONRPCError(&eth.JSONRPCError{
			Code:    100,
			Message: err1.Error(),
		})

		return
	}

	log.Println("errorHandler", err.Error())

	c.Echo().DefaultHTTPErrorHandler(err, c)
}
