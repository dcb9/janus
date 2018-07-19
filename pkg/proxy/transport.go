package proxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dcb9/janus/pkg/rpc"
	"github.com/dcb9/janus/pkg/transformer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Transport struct {
	reqTransformers  map[string]transformer.RequestTransformer
	respTransformers map[string]transformer.ResponseTransformer
	logger           log.Logger
	userInfo         *url.Userinfo
}

func (t *Transport) RoundTrip(httpReq *http.Request) (resp *http.Response, err error) {
	rpcReq := &rpc.JSONRPCRequest{}
	if err = bind(httpReq, &rpcReq); err != nil {
		return nil, err
	}

	defer func(id json.RawMessage) {
		if err != nil {
			switch err.(type) {
			case *rpc.JSONRPCError:
				resp, err = newJSONResponse(http.StatusInternalServerError, &rpc.JSONRPCResult{
					Error: err.(*rpc.JSONRPCError),
					ID:    id,
				})
			}
		}
	}(rpcReq.ID)

	if trafo, ok := t.reqTransformers[rpcReq.Method]; ok {
		if rpcReq, err = trafo.Transform(rpcReq); err != nil {
			return nil, err
		}
	}

	return t.do(httpReq, rpcReq)
}

func (t *Transport) do(r *http.Request, bodyI interface{}) (*http.Response, error) {
	n, err := t.deriveReq(r, bodyI)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultTransport.RoundTrip(n)
	if err != nil {
		return nil, err
	}

	var respBodyBytes []byte
	resp.Body, respBodyBytes = copyBody(resp.Body)

	level.Debug(t.logger).Log("respBody", respBodyBytes, "status", resp.Status, "statusCode", resp.StatusCode)

	return resp, nil
}

func (t *Transport) deriveReq(r *http.Request, bodyI interface{}) (*http.Request, error) {
	n := *r
	body, err := json.Marshal(bodyI)
	if err != nil {
		return nil, err
	}

	password, _ := t.userInfo.Password()
	n.SetBasicAuth(t.userInfo.Username(), password)
	n.ContentLength = int64(len(body))
	n.Body = ioutil.NopCloser(bytes.NewReader(body))

	level.Debug(t.logger).Log("url", n.URL.String(), "newRequestBody", body)

	return &n, nil
}
