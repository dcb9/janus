package qtum

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type Client struct {
	URL  string
	doer doer

	logger log.Logger
	debug  bool

	id      *big.Int
	idStep  *big.Int
	idMutex sync.Mutex
}

func NewClient(rpcURL string, opts ...func(*Client) error) (*Client, error) {
	err := checkRPCURL(rpcURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		doer:   http.DefaultClient,
		URL:    rpcURL,
		logger: log.NewNopLogger(),
		debug:  false,
		id:     big.NewInt(0),
		idStep: big.NewInt(1),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) Request(method string, params interface{}, result interface{}) error {
	r, err := c.NewRPCRequest(method, params)
	if err != nil {
		return err
	}

	resp, err := c.Do(r)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resp.RawResult, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) Do(req *JSONRPCRequest) (*SuccessJSONRPCResult, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	l := log.With(level.Debug(c.logger), "method", req.Method)
	if c.debug {
		l.Log("reqBody", reqBody)
	}

	respBody, err := c.do(bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrap(err, "Client#do")
	}

	if c.debug {
		l.Log("respBody", respBody)
	}

	res, err := responseBodyToResult(respBody)
	if err != nil {
		return nil, errors.Wrap(err, "responseBodyToResult")
	}

	return res, nil
}

func (c *Client) NewRPCRequest(method string, params interface{}) (*JSONRPCRequest, error) {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	c.idMutex.Lock()
	c.id = c.id.Add(c.id, c.idStep)
	c.idMutex.Unlock()

	return &JSONRPCRequest{
		JSONRPC: RPCVersion,
		ID:      json.RawMessage(`"` + c.id.String() + `"`),
		Method:  method,
		Params:  paramsJSON,
	}, nil
}

func (c *Client) do(body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.URL, body)
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

	return ioutil.ReadAll(resp.Body)
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
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

func responseBodyToResult(body []byte) (*SuccessJSONRPCResult, error) {
	var res *JSONRPCResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return &SuccessJSONRPCResult{
		ID:        res.ID,
		RawResult: res.RawResult,
		JSONRPC:   res.JSONRPC,
	}, nil
}

func checkRPCURL(u string) error {
	if u == "" {
		return errors.New("URL must be set")
	}

	qtumRPC, err := url.Parse(u)
	if err != nil {
		return errors.Errorf("QTUM_RPC URL: %s", u)
	}

	if qtumRPC.User == nil {
		return errors.Errorf("QTUM_RPC URL (must specify user & password): %s", u)
	}

	return nil
}
