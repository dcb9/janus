package server

import (
	"encoding/json"
	"net/http"

	"github.com/dcb9/janus/pkg/rpc"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log/level"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (s *Server) httpHandler(c echo.Context) error {
	var rpcReq rpc.JSONRPCRequest
	if err := c.Bind(&rpcReq); err != nil {
		return err
	}
	c.Set("rpcID", rpcReq.ID)

	level.Info(s.logger).Log("msg", "before transform request", "method", rpcReq.Method, "params", []byte(rpcReq.Params))
	responseTransformer, err := s.transformerManager.TransformRequest(&rpcReq)

	if err != nil {
		level.Error(s.logger).Log("transformErr", err)
		return err
	}

	level.Info(s.logger).Log("msg", "after transform request", "method", rpcReq.Method, "params", []byte(rpcReq.Params))

	result, err := s.qtumClient.Request(&rpcReq)
	if err != nil {
		return errors.Wrap(err, "qtum client request")
	}
	if result.RawResult == nil {
		return errors.New("server response result.RawResult is nil")
	}

	if responseTransformer != nil {
		level.Info(s.logger).Log("msg", "before transform response", "rawResult", []byte(result.RawResult))
		newResult, err := responseTransformer(result.RawResult)
		if err != nil {
			level.Error(s.logger).Log("msg", "transform response", "err", err)
			return err
		}
		if result.RawResult, err = json.Marshal(newResult); err != nil {
			return err
		}
		level.Info(s.logger).Log("msg", "after transform response", "rawResult", []byte(result.RawResult))
	}

	result.JSONRPC = "2.0"
	return c.JSON(http.StatusOK, &result)
}

func (s *Server) errorHandler(err error, c echo.Context) {
	if err != nil {
		id, _ := c.Get("rpcID").(json.RawMessage)
		err := errors.Cause(err)
		if err == transformer.UnmarshalRequestErr {
			if err := c.JSON(http.StatusInternalServerError, &rpc.JSONRPCResult{
				Error: &rpc.JSONRPCError{
					Code:    rpc.ErrInvalid,
					Message: "Input is invalid",
				},
				ID:      id,
				JSONRPC: "2.0",
			}); err != nil {
				level.Error(s.logger).Log("msg", "reply to client error", "err", err)
			}
			return
		}

		switch err.(type) {
		case *rpc.JSONRPCError:
			if err := c.JSON(http.StatusInternalServerError, &rpc.JSONRPCResult{
				Error:   err.(*rpc.JSONRPCError),
				ID:      id,
				JSONRPC: "2.0",
			}); err != nil {
				level.Error(s.logger).Log("msg", "reply to client error", "err", err)
			}
			return
		}
	}

	level.Error(s.logger).Log("err", err)
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
