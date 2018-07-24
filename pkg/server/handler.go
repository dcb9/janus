package server

import (
	"encoding/json"
	"log"

	"github.com/dcb9/janus/pkg/rpc"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log/level"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func httpHandler(c *myCtx) (interface{}, error) {
	rpcReq := c.rpcReq

	switch rpcReq.Method {
	case "personal_unlockAccount":
		return true, nil
	}

	level.Info(c.logger).Log("msg", "before transform request", "method", rpcReq.Method, "params", []byte(rpcReq.Params))
	responseTransformer, err := c.server.transformerManager.TransformRequest(rpcReq)

	if err != nil {
		level.Error(c.logger).Log("transformErr", err)
		return nil, err
	}

	level.Info(c.logger).Log("msg", "after transform request", "method", rpcReq.Method, "params", []byte(rpcReq.Params))

	result, err := c.server.qtumClient.Request(rpcReq)
	if err != nil {
		return nil, errors.Wrap(err, "qtum client request")
	}
	if result.RawResult == nil {
		return nil, errors.New("server response result.RawResult is nil")
	}

	if responseTransformer != nil {
		level.Info(c.logger).Log("msg", "before transform response", "rawResult", []byte(result.RawResult))
		newResult, err := responseTransformer(result.RawResult)
		if err != nil {
			level.Error(c.logger).Log("msg", "transform response", "err", err)
			return nil, err
		}
		if result.RawResult, err = json.Marshal(newResult); err != nil {
			return nil, err
		}
		level.Info(c.logger).Log("msg", "after transform response", "rawResult", []byte(result.RawResult))
	}

	return result.RawResult, nil
}

func errorHandler(err error, c echo.Context) {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if ok {
		err := errors.Cause(err)
		if err == transformer.UnmarshalRequestErr {
			if err := cc.JSONRPCError(&rpc.JSONRPCError{
				Code:    rpc.ErrInvalid,
				Message: "Input is invalid",
			}); err != nil {
				level.Error(cc.logger).Log("msg", "reply to client error", "err", err)
			}
			return
		}

		switch err.(type) {
		case *rpc.JSONRPCError:
			if err := cc.JSONRPCError(err.(*rpc.JSONRPCError)); err != nil {
				level.Error(cc.logger).Log("msg", "reply to client error", "err", err)
			}
			return
		}

		if err := cc.JSONRPCError(&rpc.JSONRPCError{
			Code:    rpc.ErrInvalid,
			Message: err.Error(),
		}); err != nil {
			level.Error(cc.logger).Log("msg", "reply to client error", "err", err)
		}
		return
	}

	log.Println("errorHandler", err.Error())

	c.Echo().DefaultHTTPErrorHandler(err, c)
}
