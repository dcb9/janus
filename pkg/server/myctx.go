package server

import (
	"encoding/json"
	"net/http"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log"
	"github.com/labstack/echo"
)

type myCtx struct {
	echo.Context
	rpcReq      *eth.JSONRPCRequest
	logger      log.Logger
	transformer *transformer.Transformer
}

func (c *myCtx) JSONRPCResult(result interface{}) error {
	response, err := eth.NewJSONRPCResult(c.rpcReq.ID, result)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (c *myCtx) JSONRPCError(err *eth.JSONRPCError) error {
	var id json.RawMessage
	if c.rpcReq != nil && c.rpcReq.ID != nil {
		id = c.rpcReq.ID
	}
	resp := &eth.JSONRPCResult{
		ID:      id,
		Error:   err,
		JSONRPC: eth.RPCVersion,
	}

	if !c.Response().Committed {
		return c.JSON(http.StatusInternalServerError, resp)
	}

	return nil
}
