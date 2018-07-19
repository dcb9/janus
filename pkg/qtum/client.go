package qtum

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"

	"bytes"

	"io"
	"io/ioutil"

	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type Client struct {
	rpcURL  string
	doer    doer
	logger  log.Logger
	debug   bool
	id      *big.Int
	idMutex sync.Mutex
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func NewClient(rpcURL string, opts ...func(*Client) error) (*Client, error) {
	c := &Client{
		doer:   http.DefaultClient,
		rpcURL: rpcURL,
		logger: log.NewNopLogger(),
		debug:  false,
		id:     big.NewInt(0),
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}
func SetDoer(d doer) func(*Client) error {
	return func(c *Client) error {
		c.doer = d
		return nil
	}
}

func SetDebug(debug bool) func(*Client) error {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

func SetLogger(l log.Logger) func(*Client) error {
	return func(c *Client) error {
		c.logger = log.WithPrefix(l, "component", "qtum.Client")
		return nil
	}
}

func (c *Client) GetHexAddress(addr string) (string, error) {
	r := c.NewRPCRequest(MethodGethexaddress)
	r.Params = json.RawMessage(fmt.Sprintf(`["%s"]`, addr))

	res, err := c.Request(r)
	if err != nil {
		return "", err
	}

	return string(res.RawResult), nil
}

func (c *Client) FromHexAddress(addr string) (string, error) {
	level.Debug(c.logger).Log("fromHexAddress", addr)
	r := c.NewRPCRequest(MethodFromhexaddress)
	r.Params = json.RawMessage(fmt.Sprintf(`["%s"]`, addr))
	res, err := c.Request(r)
	if err != nil {
		return "", err
	}

	return string(res.RawResult), nil
}

func (c *Client) Request(reqBody *rpc.JSONRPCRequest) (*rpc.JSONRPCResult, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	res, err := c.do(bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "Client#do")
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return res, nil
}

func (c *Client) do(body io.Reader) (*rpc.JSONRPCResult, error) {
	req, err := http.NewRequest(http.MethodPost, c.rpcURL, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res *rpc.JSONRPCResult
	if err = json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}

	return res, nil
}

var step = big.NewInt(1)

func (c *Client) NewRPCRequest(method string) *rpc.JSONRPCRequest {
	c.idMutex.Lock()
	c.id = c.id.Add(c.id, step)
	c.idMutex.Unlock()
	return &rpc.JSONRPCRequest{
		JSONRPC: Version,
		ID:      json.RawMessage(`"` + c.id.String() + `"`),
		Method:  method,
	}
}
