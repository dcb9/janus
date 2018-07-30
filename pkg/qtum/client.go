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

	"github.com/bitly/go-simplejson"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

// FIXME: Abstract all API calls into a RPC call helper. See:
// https://github.com/hayeah/eostools/blob/8e1d0d48c74b7ca6a74b132d619a4e7f8673d26a/eos-actions/main.go#L16

// FIXME: We can probably remove all methods from Client, and only provide a single `Request` method

// FIXME: Define all RPC types in rpcTypes.go

// FIXME: rename Client -> RPC or RPCClient
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

// FIXME: GetHexAddress -> Base58AddressToHex
func (c *Client) GetHexAddress(addr string) (string, error) {
	r := c.NewRPCRequest(MethodGethexaddress)
	r.Params = json.RawMessage(fmt.Sprintf(`["%s"]`, addr))

	res, err := c.Request(r)
	if err != nil {
		return "", err
	}

	var hexAddr string
	if err = json.Unmarshal(res.RawResult, &hexAddr); err != nil {
		return "", err
	}

	return hexAddr, nil
}

func (c *Client) FromHexAddress(addr string) (string, error) {
	level.Debug(c.logger).Log("fromHexAddress", addr)
	r := c.NewRPCRequest(MethodFromhexaddress)
	r.Params = json.RawMessage(fmt.Sprintf(`["%s"]`, addr))
	res, err := c.Request(r)
	if err != nil {
		return "", err
	}

	var qtumAddr string
	if err = json.Unmarshal(res.RawResult, &qtumAddr); err != nil {
		return "", err
	}

	return qtumAddr, nil
}

// FIXME: Define the types for all API methods: [methodName]Request, [methodName]Response
//
// type GetTransactionReceiptResponse []struct {
// 	// A int `json:"a"`
// }

func (c *Client) GetTransactionReceipt(txHash string) (*TransactionReceipt, error) {
	r := c.NewRPCRequest(MethodGettransactionreceipt)
	r.Params = json.RawMessage(fmt.Sprintf(`["%s"]`, txHash))

	result, err := c.Request(r)
	if err != nil {
		return nil, err
	}

	js, err := simplejson.NewJson(result.RawResult)
	if err != nil {
		return nil, err
	}
	receiptJSON, err := js.GetIndex(0).Encode()
	if err != nil {
		return nil, err
	}
	var receipt *TransactionReceipt
	err = json.Unmarshal(receiptJSON, &receipt)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (c *Client) DecodeRawTransaction(hex string) (*DecodedRawTransaction, error) {
	r := c.NewRPCRequest(MethodDecoderawtransaction)
	r.Params = json.RawMessage(fmt.Sprintf(`["%s"]`, hex))

	result, err := c.Request(r)
	if err != nil {
		return nil, err
	}

	var tx *DecodedRawTransaction
	if err = json.Unmarshal(result.RawResult, &tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *Client) Request(reqBody *rpc.JSONRPCRequest) (*rpc.SuccessJSONRPCResult, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	respBody, err := c.do(bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "Client#do")
	}
	res, err := responseBodyToResult(respBody)
	if err != nil {
		return nil, errors.Wrap(err, "responseBodyToResult")
	}

	return res, nil
}

func (c *Client) do(body io.Reader) ([]byte, error) {
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

	return ioutil.ReadAll(resp.Body)
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

func responseBodyToResult(body []byte) (*rpc.SuccessJSONRPCResult, error) {
	var res *rpc.JSONRPCResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return &rpc.SuccessJSONRPCResult{
		ID:        res.ID,
		RawResult: res.RawResult,
		JSONRPC:   res.JSONRPC,
	}, nil
}
