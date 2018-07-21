package qtum

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	Version = "1.0"
)

const (
	MethodGethexaddress         = "gethexaddress"
	MethodFromhexaddress        = "fromhexaddress"
	MethodSendtocontract        = "sendtocontract"
	MethodGettransactionreceipt = "gettransactionreceipt"
	MethodGettransaction        = "gettransaction"
	MethodCreatecontract        = "createcontract"
	MethodSendtoaddress         = "sendtoaddress"
	MethodCallcontract          = "callcontract"
	MethodDecoderawtransaction  = "decoderawtransaction"
)

type (
	TransactionReceipt struct {
		BlockHash         string `json:"blockHash"`
		BlockNumber       uint64 `json:"blockNumber"`
		TransactionHash   string `json:"transactionHash"`
		TransactionIndex  uint64 `json:"transactionIndex"`
		From              string `json:"from"`
		To                string `json:"to"`
		CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`
		GasUsed           uint64 `json:"gasUsed"`
		ContractAddress   string `json:"contractAddress"`
		Excepted          string `json:"excepted"`
		Log               []Log  `json:"log"`
	}

	Log struct {
		Address string   `json:"address"`
		Topics  []string `json:"topics"`
		Data    string   `json:"data"`
	}

	TransactionInV struct {
		Txid      string `json:"txid"`
		Vout      int64  `json:"vout"`
		ScriptSig struct {
			Asm string `json:"asm"`
			Hex string `json:"hex"`
		} `json:"scriptSig"`
		Txinwitness []string `json:"txinwitness"`
		Sequence    int64    `json:"sequence"`
	}

	TransactionOutV struct {
		Value        float64 `json:"value"`
		N            int64   `json:"n"`
		ScriptPubKey struct {
			Asm       string   `json:"asm"`
			Hex       string   `json:"hex"`
			ReqSigs   int64    `json:"reqSigs"`
			Type      string   `json:"type"`
			Addresses []string `json:"addresses"`
		} `json:"scriptPubKey"`
	}

	Transaction struct {
		Txid     string             `json:"txid"`
		Hash     string             `json:"hash"`
		Size     int64              `json:"size"`
		Vsize    int64              `json:"vsize"`
		Version  int64              `json:"version"`
		Locktime int64              `json:"locktime"`
		Vin      []*TransactionInV  `json:"vin"`
		Vout     []*TransactionOutV `json:"vout"`
	}

	ASM struct {
		VMVersion  string
		GasLimit   string
		GasPrice   string
		Instructor string
	}
	CallASM struct {
		ASM
		EncodedABI      string
		ContractAddress string
	}

	CreateASM struct {
		ASM
		EncodedABI string
	}
)

func ParseCallASM(asm string) (*CallASM, error) {
	parts := strings.Split(asm, " ")
	if len(parts) < 6 {
		return nil, errors.New("invalid call sam")
	}

	return &CallASM{
		ASM: ASM{
			VMVersion:  parts[0],
			GasLimit:   parts[1],
			GasPrice:   parts[2],
			Instructor: parts[5],
		},
		EncodedABI:      parts[3],
		ContractAddress: parts[4],
	}, nil
}

func ParseCreateASM(asm string) (*CreateASM, error) {
	parts := strings.Split(asm, " ")
	if len(parts) < 5 {
		return nil, errors.New("invalid create sam")
	}

	return &CreateASM{
		ASM: ASM{
			VMVersion:  parts[0],
			GasLimit:   parts[1],
			GasPrice:   parts[2],
			Instructor: parts[4],
		},
		EncodedABI: parts[3],
	}, nil
}

func (asm *ASM) GetGasPrice() (*big.Int, error) {
	return stringToBigInt(asm.GasPrice)
}

func (asm *ASM) GetGasLimit() (*big.Int, error) {
	return stringToBigInt(asm.GasLimit)
}

func (asm *CreateASM) GetEncodedABI() string {
	return asm.EncodedABI
}

func (asm *CallASM) GetEncodedABI() string {
	return asm.EncodedABI
}

func stringToBigInt(str string) (*big.Int, error) {
	var success bool
	v := new(big.Int)
	if v, success = v.SetString(str, 10); !success {
		return nil, errors.New(fmt.Sprintf("failed to parse str: %s to big.Int", str))
	}
	return v, nil
}
