package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dcb9/janus/pkg/qtum"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type Transport struct {
	transformers map[string]transformer
	logger       log.Logger
	userInfo     *url.Userinfo
}

func (t *Transport) RoundTrip(httpReq *http.Request) (*http.Response, error) {
	rpcReq := &qtum.JSONRPCRequest{}
	bodyBytes, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		return nil, err
	}
	httpReq.Body.Close()

	level.Debug(t.logger).Log("rawRequestBody", bodyBytes, "len", len(bodyBytes))
	if err := json.Unmarshal(bodyBytes, rpcReq); err != nil {
		return nil, err
	}

	newHTTPReq := *httpReq

	if transformer, ok := t.transformers[rpcReq.Method]; ok {
		if rpcReq, err = transformer.transform(rpcReq); err != nil {
			return nil, errors.Wrap(err, "transform")
		}
	}

	newBodyBytes, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	password, _ := t.userInfo.Password()
	newHTTPReq.SetBasicAuth(t.userInfo.Username(), password)
	newHTTPReq.ContentLength = int64(len(newBodyBytes))
	newHTTPReq.Body = ioutil.NopCloser(bytes.NewReader(newBodyBytes))

	level.Debug(t.logger).Log("url", newHTTPReq.URL.String(), "newRequestBody", newBodyBytes, "len", len(newBodyBytes))

	resp, err := http.DefaultTransport.RoundTrip(&newHTTPReq)
	if err != nil {
		return nil, err
	}

	respBodyBytes := []byte{}
	resp.Body, respBodyBytes = copyBody(resp.Body)

	level.Debug(t.logger).Log("respBody", respBodyBytes, "status", resp.Status, "statusCode", resp.StatusCode)

	return resp, nil
}

func copyBody(r io.ReadCloser) (io.ReadCloser, []byte) {
	body := make([]byte, 0)
	if r != nil {
		body, _ = ioutil.ReadAll(r)
		r = ioutil.NopCloser(bytes.NewReader(body))
	}
	return r, body
}
