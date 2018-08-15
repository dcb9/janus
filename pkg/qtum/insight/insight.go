package insight

import (
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/pkg/errors"
)

const (
	MainNetEndpoint = "https://explorer.qtum.org/insight-api"
	TestNetEndpoint = "https://testnet.qtum.org/insight-api"
)

type Insight struct {
	Endpoint string
}

func New(endpoint string) *Insight {
	endpoint = strings.TrimRight(endpoint, "/")

	return &Insight{
		Endpoint: endpoint,
	}
}

func (i *Insight) GetAddressSummary(addr string) (*AddressSummary, error) {
	var summary AddressSummary
	if err := i.request(fmt.Sprintf("/addr/%s", addr), &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func (i *Insight) GetBlock(blockHash string) (*Block, error) {
	var block Block
	if err := i.request(fmt.Sprintf("/block/%s", &blockHash), &block); err != nil {
		return nil, err
	}

	return &block, nil
}

func (i *Insight) GetRawBlock(blockHash string) (*RawBlock, error) {
	var rawBlock RawBlock
	if err := i.request(fmt.Sprintf("/rawblock/%s", blockHash), &rawBlock); err != nil {
		return nil, err
	}

	return &rawBlock, nil
}

func (i *Insight) GetStatus() (*Status, error) {
	var status Status
	if err := i.request("/status", &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (i *Insight) GetBlockIndex(height int64) (*BlockIndex, error) {
	var blockIndex BlockIndex
	if err := i.request(fmt.Sprintf("/block-index/%d", height), &blockIndex); err != nil {
		return nil, err
	}

	return &blockIndex, nil
}

func (i *Insight) request(url string, body interface{}) error {
	if i.Endpoint == "" {
		return errors.New("insight api endpoint cannot be empty (regtest network is not supported yet)")
	}

	url = i.Endpoint + url
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return i.rawRequest(req, body)
}

func (i *Insight) rawRequest(req *http.Request, body interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return errors.New("body.Body is nil")
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(body)
}

/*
{
  "hash": "0000de7d93ecdfc31efa7c66d10b87c41bffee5d71c6df7dd832463a3ccb17fc",
  "size": 327,
  "height": 1,
  "version": 536870912,
  "merkleroot": "dee4bb3300265f05bb5f179eb9c9124cf9a0ff596b55d1af56e8032e435a5740",
  "tx": [
    "dee4bb3300265f05bb5f179eb9c9124cf9a0ff596b55d1af56e8032e435a5740"
  ],
  "time": 1506743769,
  "nonce": 44182,
  "bits": "1f00ffff",
  "difficulty": 0.1525,
  "chainwork": "0000000000000000000000000000000000000000000000000000000000020002",
  "confirmations": 194655,
  "previousblockhash": "0000e803ee215c0684ca0d2f9220594d3f828617972aad66feb2ba51f5e14222",
  "nextblockhash": "0000715144ff7be935a58304f68dd5c09514c0a0291444d72c1d4e09c8163249",
  "flags": "proof-of-work",
  "reward": 20000,
  "isMainChain": true,
  "minedBy": "qSmstCdKdqoRcSZueiRZ3KtL2xJRWsUeFZ",
  "poolInfo": {}
}
*/
type Block struct {
	Hash              string   `json:"hash"`
	Size              int64    `json:"size"`
	Height            int64    `json:"height"`
	Version           int64    `json:"version"`
	Merkleroot        string   `json:"merkleroot"`
	Tx                []string `json:"tx"`
	Time              int64    `json:"time"`
	Nonce             int64    `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	Confirmations     int64    `json:"confirmations"`
	Previousblockhash string   `json:"previousblockhash"`
	Nextblockhash     string   `json:"nextblockhash"`
	Flags             string   `json:"flags"`
	Reward            int64    `json:"reward"`
	IsMainChain       bool     `json:"isMainChain"`
	MinedBy           string   `json:"minedBy"`
	PoolInfo          struct{} `json:"poolInfo"`
}

/**
{"rawblock":"000000202242e1f551b..."}
*/
type RawBlock struct {
	RawBlock string `json:"rawblock"`
}

/**
{"blockHash":"0000de7d93ecdfc31efa7c66d10b87c41bffee5d71c6df7dd832463a3ccb17fc"}
*/
type BlockIndex struct {
	BlockHash string `json:"blockHash"`
}

/*
{
  "addrStr": "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
  "balance": 90.20163883,
  "balanceSat": 9020163883,
  "totalReceived": 638.16579884,
  "totalReceivedSat": 63816579884,
  "totalSent": 547.96416001,
  "totalSentSat": 54796416001,
  "unconfirmedBalance": 0,
  "unconfirmedBalanceSat": 0,
  "unconfirmedTxApperances": 0,
  "txApperances": 42,
  "transactions": [
    "508c63cbcdc9295e3771b3411cce86cc6b2719b727a9a87c437e3ff9df112231"
  ]
}
*/
type AddressSummary struct {
	AddrStr                 string   `json:"addrStr"`
	Balance                 float64  `json:"balance"`
	BalanceSat              int64    `json:"balanceSat"`
	TotalReceived           float64  `json:"totalReceived"`
	TotalReceivedSat        int64    `json:"totalReceivedSat"`
	TotalSent               float64  `json:"totalSent"`
	TotalSentSat            int64    `json:"totalSentSat"`
	UnconfirmedBalance      float64  `json:"unconfirmedBalance"`
	UnconfirmedBalanceSat   int64    `json:"unconfirmedBalanceSat"`
	UnconfirmedTxApperances int64    `json:"unconfirmedTxApperances"`
	TxApperances            int64    `json:"txApperances"`
	Transactions            []string `json:"transactions"`
}

/*
{
  "info": {
    "version": 141600,
    "protocolversion": 70016,
    "walletversion": 130000,
    "balance": 2559872.194715,
    "blocks": 789,
    "timeoffset": 0,
    "connections": 0,
    "proxy": "",
    "difficulty": {
      "proof-of-work": 4.656542373906925e-10,
      "proof-of-stake": 4.656542373906925e-10
    },
    "testnet": false,
    "keypoololdest": 1534245568,
    "keypoolsize": 100,
    "paytxfee": 0,
    "relayfee": 0.004,
    "errors": "",
    "network": "livenet",
    "reward": 2000000000000
  }
}
*/
type Status struct {
	Info struct {
		Version         int64   `json:"version"`
		Protocolversion int64   `json:"protocolversion"`
		Walletversion   int64   `json:"walletversion"`
		Balance         float64 `json:"balance"`
		Blocks          int64   `json:"blocks"`
		Timeoffset      int64   `json:"timeoffset"`
		Connections     int64   `json:"connections"`
		Proxy           string  `json:"proxy"`
		Difficulty      struct {
			ProofOfWork  float64 `json:"proof-of-work"`
			ProofOfStake float64 `json:"proof-of-stake"`
		} `json:"difficulty"`
		Testnet       bool    `json:"testnet"`
		Keypoololdest int64   `json:"keypoololdest"`
		Keypoolsize   int64   `json:"keypoolsize"`
		Paytxfee      int64   `json:"paytxfee"`
		Relayfee      float64 `json:"relayfee"`
		Errors        string  `json:"errors"`
		Network       string  `json:"network"`
		Reward        int64   `json:"reward"`
	} `json:"info"`
}
