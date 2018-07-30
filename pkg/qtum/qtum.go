package qtum

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

const (
	Version = "1.0"
)

const (
	// FIXME: camel case. MethodGethexaddress -> MethodGetHexAddress
	MethodGethexaddress         = "gethexaddress"
	MethodFromhexaddress        = "fromhexaddress"
	MethodSendtocontract        = "sendtocontract"
	MethodGettransactionreceipt = "gettransactionreceipt"
	MethodGettransaction        = "gettransaction"
	MethodCreatecontract        = "createcontract"
	MethodSendtoaddress         = "sendtoaddress"
	MethodCallcontract          = "callcontract"
	MethodDecoderawtransaction  = "decoderawtransaction"
	MethodGetblockcount         = "getblockcount"
	MethodGetblockchaininfo     = "getblockchaininfo"
	MethodSearchlogs            = "searchlogs"
	MethodWaitforlogs           = "waitforlogs"
)

// FIXME: rename TransactionReceipt -> GetTransactionReceiptResponse
// FIXME: rename all other qtum RPC types
// FIXME: move all RPC types into its own file.
// FIXME: (optional) add RPC examples to each struct
type (
	/* example:
	{
		"blockhash": "afafafa..."
	}
	*/
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

	DecodedRawTransactionInV struct {
		Txid      string `json:"txid"`
		Vout      int64  `json:"vout"`
		ScriptSig struct {
			Asm string `json:"asm"`
			Hex string `json:"hex"`
		} `json:"scriptSig"`
		Txinwitness []string `json:"txinwitness"`
		Sequence    int64    `json:"sequence"`
	}

	DecodedRawTransactionOutV struct {
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

	DecodedRawTransaction struct {
		Txid     string                       `json:"txid"`
		Hash     string                       `json:"hash"`
		Size     int64                        `json:"size"`
		Vsize    int64                        `json:"vsize"`
		Version  int64                        `json:"version"`
		Locktime int64                        `json:"locktime"`
		Vin      []*DecodedRawTransactionInV  `json:"vin"`
		Vout     []*DecodedRawTransactionOutV `json:"vout"`
	}

	TransactionDetail struct {
		Account   string  `json:"account"`
		Address   string  `json:"address"`
		Category  string  `json:"category"`
		Amount    float64 `json:"amount"`
		Label     string  `json:"label"`
		Vout      int64   `json:"vout"`
		Fee       float64 `json:"fee"`
		Abandoned bool    `json:"abandoned"`
	}

	Transaction struct {
		Amount            float64              `json:"amount"`
		Fee               float64              `json:"fee"`
		Confirmations     int64                `json:"confirmations"`
		Blockhash         string               `json:"blockhash"`
		Blockindex        int64                `json:"blockindex"`
		Blocktime         int64                `json:"blocktime"`
		Txid              string               `json:"txid"`
		Time              int64                `json:"time"`
		Timereceived      int64                `json:"timereceived"`
		Bip125Replaceable string               `json:"bip125-replaceable"`
		Details           []*TransactionDetail `json:"details"`
		Hex               string               `json:"hex"`
	}

	// FIXME: extract asm to its own file. btcasm.go

	// ASM is Bitcoin Script extended by qtum to support smart contracts
	ASM struct {
		VMVersion  string
		GasLimit   string
		GasPrice   string
		Instructor string
	}
	CallASM struct {
		ASM
		// FIXME: EncodedABI -> CallData
		EncodedABI      string
		ContractAddress string
	}
	CreateASM struct {
		ASM
		// FIXME: EncodedABI -> CallData
		EncodedABI string
	}

	BlockChainInfo struct {
		Bestblockhash string `json:"bestblockhash"`
		Bip9Softforks struct {
			Csv struct {
				Since     int64  `json:"since"`
				StartTime int64  `json:"startTime"`
				Status    string `json:"status"`
				Timeout   int64  `json:"timeout"`
			} `json:"csv"`
			Segwit struct {
				Since     int64  `json:"since"`
				StartTime int64  `json:"startTime"`
				Status    string `json:"status"`
				Timeout   int64  `json:"timeout"`
			} `json:"segwit"`
		} `json:"bip9_softforks"`
		Blocks     int64   `json:"blocks"`
		Chain      string  `json:"chain"`
		Chainwork  string  `json:"chainwork"`
		Difficulty float64 `json:"difficulty"`
		Headers    int64   `json:"headers"`
		Mediantime int64   `json:"mediantime"`
		Pruned     bool    `json:"pruned"`
		Softforks  []struct {
			ID     string `json:"id"`
			Reject struct {
				Status bool `json:"status"`
			} `json:"reject"`
			Version int64 `json:"version"`
		} `json:"softforks"`
		Verificationprogress int64 `json:"verificationprogress"`
	}
)

func ParseCallASM(asm string) (*CallASM, error) {
	parts := strings.Split(asm, " ")
	if len(parts) < 6 {
		// FIXME: typo: sam -> ASM
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
		// FIXME: typo: sam -> ASM
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

// FIXME: rename GetGasPrice -> GasPrice
func (asm *ASM) GetGasPrice() (*big.Int, error) {
	return stringToBigInt(asm.GasPrice)
}

// FIXME: rename GetGasLimit -> GasLimit
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
		// FIXME: use errors.Errorf
		return nil, errors.New(fmt.Sprintf("Failed to parse big.Int: %s", str))
	}
	return v, nil
}
