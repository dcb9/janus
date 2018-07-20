package proxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log/level"
)

func (p *proxy) RoundTrip(httpReq *http.Request) (resp *http.Response, err error) {
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

	responseTransformer, err := p.transformerManager.TransformRequest(rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = p.do(httpReq, rpcReq)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	result := new(rpc.JSONRPCResult)
	if err = json.Unmarshal(body, &result); err != nil {
		level.Error(p.logger).Log("err", err, "body", body)
		return nil, err
	}

	if p.debug {
		level.Debug(p.logger).Log("msg", "before call response transformer")
	}
	if err = responseTransformer(result); err != nil {
		return nil, err
	}
	if p.debug {
		level.Debug(p.logger).Log("msg", "after call response transformer")
	}

	newBodyBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if p.debug {
		level.Debug(p.logger).Log("newRespBody", newBodyBytes)
	}

	resp.ContentLength = int64(len(newBodyBytes))
	resp.Body = ioutil.NopCloser(bytes.NewReader(newBodyBytes))

	return resp, nil
}

func (p *proxy) do(r *http.Request, bodyI interface{}) (*http.Response, error) {
	n, err := p.deriveReq(r, bodyI)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultTransport.RoundTrip(n)
	if err != nil {
		return nil, err
	}

	if p.debug {
		var respBodyBytes []byte
		resp.Body, respBodyBytes = copyBody(resp.Body)

		level.Debug(p.logger).Log("respBody", respBodyBytes, "status", resp.Status, "statusCode", resp.StatusCode)
	}

	return resp, nil
}

func (p *proxy) deriveReq(r *http.Request, bodyI interface{}) (*http.Request, error) {
	n := *r
	body, err := json.Marshal(bodyI)
	if err != nil {
		return nil, err
	}

	password, _ := p.qtumRPC.User.Password()
	n.SetBasicAuth(p.qtumRPC.User.Username(), password)
	n.ContentLength = int64(len(body))
	n.Body = ioutil.NopCloser(bytes.NewReader(body))

	level.Debug(p.logger).Log("url", n.URL.String(), "newRequestBody", body)

	return &n, nil
}
