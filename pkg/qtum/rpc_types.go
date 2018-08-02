package qtum

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/dcb9/janus/pkg/utils"
)

type (
	Log struct {
		Address string   `json:"address"`
		Topics  []string `json:"topics"`
		Data    string   `json:"data"`
	}

	/*
		{
		  "chain": "regtest",
		  "blocks": 4137,
		  "headers": 4137,
		  "bestblockhash": "3863e43665ab15af1167df2f30a1c6f658c64704a3a2903bb0c5afde7e5d54ff",
		  "difficulty": 4.656542373906925e-10,
		  "mediantime": 1533096368,
		  "verificationprogress": 1,
		  "chainwork": "0000000000000000000000000000000000000000000000000000000000002054",
		  "pruned": false,
		  "softforks": [
		    {
		      "id": "bip34",
		      "version": 2,
		      "reject": {
		        "status": true
		      }
		    },
		    {
		      "id": "bip66",
		      "version": 3,
		      "reject": {
		        "status": true
		      }
		    },
		    {
		      "id": "bip65",
		      "version": 4,
		      "reject": {
		        "status": true
		      }
		    }
		  ],
		  "bip9_softforks": {
		    "csv": {
		      "status": "active",
		      "startTime": 0,
		      "timeout": 999999999999,
		      "since": 432
		    },
		    "segwit": {
		      "status": "active",
		      "startTime": 0,
		      "timeout": 999999999999,
		      "since": 432
		    }
		  }
		}
	*/
	GetBlockChainInfoResponse struct {
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

// ========== SendToAddress ============= //

type (
	SendToAddressRequest struct {
		Address       string
		Amount        float64
		SenderAddress string
	}
	SendToAddressResponse string
)

func (r *SendToAddressRequest) MarshalJSON() ([]byte, error) {
	/*
		1. "address"            (string, required) The qtum address to send to.
		2. "amount"             (numeric or string, required) The amount in QTUM to send. eg 0.1
		3. "comment"            (string, optional) A comment used to store what the transaction is for.
		                             This is not part of the transaction, just kept in your wallet.
		4. "comment_to"         (string, optional) A comment to store the name of the person or organization
		                             to which you're sending the transaction. This is not part of the
		                             transaction, just kept in your wallet.
		5. subtractfeefromamount  (boolean, optional, default=false) The fee will be deducted from the amount being sent.
		                             The recipient will receive less qtums than you enter in the amount field.
		6. replaceable            (boolean, optional) Allow this transaction to be replaced by a transaction with higher fees via BIP 125
		7. conf_target            (numeric, optional) Confirmation target (in blocks)
		8. "estimate_mode"      (string, optional, default=UNSET) The fee estimate mode, must be one of:
		       "UNSET"
		       "ECONOMICAL"
		       "CONSERVATIVE"
		9. "senderaddress"      (string, optional) The quantum address that will be used to send money from.
		10."changeToSender"     (bool, optional, default=false) Return the change to the sender.
	*/
	return json.Marshal([]interface{}{
		r.Address,
		r.Amount,
		"", // comment
		"", // comment_to
		false,
		nil,
		nil,
		nil,
		r.SenderAddress,
		true,
	})
}

// ========== SendToContract ============= //

type (
	SendToContractRequest struct {
		ContractAddress string
		Datahex         string
		Amount          float64
		GasLimit        *big.Int
		GasPrice        string
		SenderAddress   string
	}
	/*
		{
		  "txid": "6b7f70d8520e1ec87ba7f1ee559b491cc3028b77ae166e789be882b5d370eac9",
		  "sender": "qTKrsHUrzutdCVu3qi3iV1upzB2QpuRsRb",
		  "hash160": "6b22910b1e302cf74803ffd1691c2ecb858d3712"
		}
	*/
	SendToContractResponse struct {
		Txid    string `json:"txid"`
		Sender  string `json:"sender"`
		Hash160 string `json:"hash160"`
	}
)

func (r *SendToContractRequest) MarshalJSON() ([]byte, error) {
	/*
	   1. "contractaddress" (string, required) The contract address that will receive the funds and data.
	   2. "datahex"  (string, required) data to send.
	   3. "amount"      (numeric or string, optional) The amount in QTUM to send. eg 0.1, default: 0
	   4. gasLimit  (numeric or string, optional) gasLimit, default: 250000, max: 40000000
	   5. gasPrice  (numeric or string, optional) gasPrice Qtum price per gas unit, default: 0.0000004, min:0.0000004
	   6. "senderaddress" (string, optional) The quantum address that will be used as sender.
	   7. "broadcast" (bool, optional, default=true) Whether to broadcast the transaction or not.
	   8. "changeToSender" (bool, optional, default=true) Return the change to the sender.
	*/

	return json.Marshal([]interface{}{
		r.ContractAddress,
		r.Datahex,
		r.Amount,
		r.GasLimit,
		r.GasPrice,
		r.SenderAddress,
	})
}

// ========== CreateContract ============= //

type (
	CreateContractRequest struct {
		ByteCode      string
		GasLimit      *big.Int
		GasPrice      string
		SenderAddress string
	}
	/*
	   {
	   "txid": "d0fe0caa1b798c36da37e9118a06a7d151632d670b82d1c7dc3985577a71880f",
	   "sender": "qTKrsHUrzutdCVu3qi3iV1upzB2QpuRsRb",
	   "hash160": "6b22910b1e302cf74803ffd1691c2ecb858d3712",
	   "address": "c89a5d225f578d84a94741490c1b40889b4f7a00"
	   }
	*/
	CreateContractResponse struct {
		Txid    string `json:"txid"`
		Sender  string `json:"sender"`
		Hash160 string `json:"hash160"`
		Address string `json:"address"`
	}
)

func (r *CreateContractRequest) MarshalJSON() ([]byte, error) {
	/*
		1. "bytecode"  (string, required) contract bytcode.
		2. gasLimit  (numeric or string, optional) gasLimit, default: 2500000, max: 40000000
		3. gasPrice  (numeric or string, optional) gasPrice QTUM price per gas unit, default: 0.0000004, min:0.0000004
		4. "senderaddress" (string, optional) The quantum address that will be used to create the contract.
		5. "broadcast" (bool, optional, default=true) Whether to broadcast the transaction or not.
		6. "changeToSender" (bool, optional, default=true) Return the change to the sender.
	*/
	return json.Marshal([]interface{}{
		r.ByteCode,
		r.GasLimit,
		r.GasPrice,
		r.SenderAddress,
	})
}

// ========== CallContract ============= //

type (
	CallContractRequest struct {
		From     string
		To       string
		Data     string
		GasLimit *big.Int
	}

	/*
		{
		  "address": "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		  "executionResult": {
		    "gasUsed": 21678,
		    "excepted": "None",
		    "newAddress": "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		    "output": "0000000000000000000000000000000000000000000000000000000000000001",
		    "codeDeposit": 0,
		    "gasRefunded": 0,
		    "depositSize": 0,
		    "gasForDeposit": 0
		  },
		  "transactionReceipt": {
		    "stateRoot": "d44fc5ad43bae52f01ff7eb4a7bba904ee52aea6c41f337aa29754e57c73fba6",
		    "gasUsed": 21678,
		    "bloom": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		    "log": []
		  }
		}
	*/
	CallContractResponse struct {
		Address         string `json:"address"`
		ExecutionResult struct {
			GasUsed       int    `json:"gasUsed"`
			Excepted      string `json:"excepted"`
			NewAddress    string `json:"newAddress"`
			Output        string `json:"output"`
			CodeDeposit   int    `json:"codeDeposit"`
			GasRefunded   int    `json:"gasRefunded"`
			DepositSize   int    `json:"depositSize"`
			GasForDeposit int    `json:"gasForDeposit"`
		} `json:"executionResult"`
		TransactionReceipt struct {
			StateRoot string        `json:"stateRoot"`
			GasUsed   int           `json:"gasUsed"`
			Bloom     string        `json:"bloom"`
			Log       []interface{} `json:"log"`
		} `json:"transactionReceipt"`
	}
)

func (r *CallContractRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		utils.RemoveHexPrefix(r.To),
		utils.RemoveHexPrefix(r.Data),
		r.From,
		r.GasLimit,
	})
}

// ========== FromHexAddress ============= //

type (
	FromHexAddressRequest  string
	FromHexAddressResponse string
)

func (r FromHexAddressRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		string(r),
	})
}

// ========== GetHexAddress ============= //

type (
	GetHexAddressRequest  string
	GetHexAddressResponse string
)

func (r GetHexAddressRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		string(r),
	})
}

// ========== DecodeRawTransaction ============= //
func (r DecodeRawTransactionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		string(r),
	})
}

type (
	DecodeRawTransactionRequest string

	/*
		{
		  "txid": "d0fe0caa1b798c36da37e9118a06a7d151632d670b82d1c7dc3985577a71880f",
		  "hash": "d0fe0caa1b798c36da37e9118a06a7d151632d670b82d1c7dc3985577a71880f",
		  "version": 2,
		  "size": 552,
		  "vsize": 552,
		  "locktime": 608,
		  "vin": [
		    {
		      "txid": "7f5350dc474f2953a3f30282c1afcad2fb61cdcea5bd949c808ecc6f64ce1503",
		      "vout": 0,
		      "scriptSig": {
		        "asm": "3045022100af4de764705dbd3c0c116d73fe0a2b78c3fab6822096ba2907cfdae2bb28784102206304340a6d260b364ef86d6b19f2b75c5e55b89fb2f93ea72c05e09ee037f60b[ALL] 03520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140",
		        "hex": "483045022100af4de764705dbd3c0c116d73fe0a2b78c3fab6822096ba2907cfdae2bb28784102206304340a6d260b364ef86d6b19f2b75c5e55b89fb2f93ea72c05e09ee037f60b012103520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140"
		      },
		      "sequence": 4294967294
		    }
		  ],
		  "vout": [
		    {
		      "value": 0,
		      "n": 0,
		      "scriptPubKey": {
		        "asm": "4 2500000 40 608060405234801561001057600080fd5b50604051602080610131833981016040525160005560fe806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b50607660cc565b60408051918252519081900360200190f35b600054604080513381526020810192909252805183927f61ec51fdd1350b55fc6e153e60509e993f8dcb537fe4318c45a573243d96cab492908290030190a2600055565b600054905600a165627a7a723058200541c7c0da642ef9004daeb68d281a3c2341e765336f10b4a0ab45dbb7b7f83c00290000000000000000000000000000000000000000000000000000000000000064 OP_CREATE",
		        "hex": "010403a0252601284d5101608060405234801561001057600080fd5b50604051602080610131833981016040525160005560fe806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b50607660cc565b60408051918252519081900360200190f35b600054604080513381526020810192909252805183927f61ec51fdd1350b55fc6e153e60509e993f8dcb537fe4318c45a573243d96cab492908290030190a2600055565b600054905600a165627a7a723058200541c7c0da642ef9004daeb68d281a3c2341e765336f10b4a0ab45dbb7b7f83c00290000000000000000000000000000000000000000000000000000000000000064c1",
		        "type": "create"
		      }
		    },
		    {
		      "value": 19996.59434,
		      "n": 1,
		      "scriptPubKey": {
		        "asm": "OP_DUP OP_HASH160 ce7137386121f7531f716d2d4ff36805bc65b3ec OP_EQUALVERIFY OP_CHECKSIG",
		        "hex": "76a914ce7137386121f7531f716d2d4ff36805bc65b3ec88ac",
		        "reqSigs": 1,
		        "type": "pubkeyhash",
		        "addresses": [
		          "qcNwyuvvPhiN4JVgwPp4QWPiK1p7YGvkf1"
		        ]
		      }
		    }
		  ]
		}
	*/
	DecodedRawTransactionResponse struct {
		Txid     string                       `json:"txid"`
		Hash     string                       `json:"hash"`
		Size     int64                        `json:"size"`
		Vsize    int64                        `json:"vsize"`
		Version  int64                        `json:"version"`
		Locktime int64                        `json:"locktime"`
		Vin      []*DecodedRawTransactionInV  `json:"vin"`
		Vout     []*DecodedRawTransactionOutV `json:"vout"`
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
)

// ========== GetTransactionReceipt ============= //
type (
	GetTransactionReceiptRequest  string
	GetTransactionReceiptResponse TransactionReceiptStruct
	/*
	   {
	     "blockHash": "975326b65c20d0b8500f00a59f76b08a98513fff7ce0484382534a47b55f8985",
	     "blockNumber": 4063,
	     "transactionHash": "c1816e5fbdd4d1cc62394be83c7c7130ccd2aadefcd91e789c1a0b33ec093fef",
	     "transactionIndex": 2,
	     "from": "6b22910b1e302cf74803ffd1691c2ecb858d3712",
	     "to": "db46f738bf32cdafb9a4a70eb8b44c76646bcaf0",
	     "cumulativeGasUsed": 68572,
	     "gasUsed": 68572,
	     "contractAddress": "db46f738bf32cdafb9a4a70eb8b44c76646bcaf0",
	     "excepted": "None",
	     "log": [
	       {
	         "address": "db46f738bf32cdafb9a4a70eb8b44c76646bcaf0",
	         "topics": [
	           "0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885",
	           "0000000000000000000000006b22910b1e302cf74803ffd1691c2ecb858d3712"
	         ],
	         "data": "0000000000000000000000000000000000000000000000000000000000000001"
	       },
	       {
	         "address": "db46f738bf32cdafb9a4a70eb8b44c76646bcaf0",
	         "topics": [
	           "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
	           "0000000000000000000000000000000000000000000000000000000000000000",
	           "0000000000000000000000006b22910b1e302cf74803ffd1691c2ecb858d3712"
	         ],
	         "data": "0000000000000000000000000000000000000000000000000000000000000001"
	       }
	     ]
	   }
	*/
	TransactionReceiptStruct struct {
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
)

func (r GetTransactionReceiptRequest) MarshalJSON() ([]byte, error) {
	/*
		1. "hash"          (string, required) The transaction hash
	*/
	return json.Marshal([]interface{}{
		string(r),
	})
}

var EmptyResponseErr = errors.New("result is empty")

func (r *GetTransactionReceiptResponse) UnmarshalJSON(data []byte) error {
	type Response GetTransactionReceiptResponse
	var resp []Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	if len(resp) == 0 {
		return EmptyResponseErr
	}

	*r = GetTransactionReceiptResponse(resp[0])

	return nil
}

// ========== GetBlockCount ============= //

type (
	GetBlockCountResponse big.Int
)

func (r *GetBlockCountResponse) UnmarshalJSON(data []byte) error {
	var i big.Int
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	*r = GetBlockCountResponse(i)
	return nil
}

// ========== GetTransaction ============= //

type (
	GetTransactionRequest struct {
		Txid string
	}

	/*
		{
		    "amount": 0,
		    "fee": -0.2012,
		    "confirmations": 2,
		    "blockhash": "ea26fd59a2145dcecd0e2f81b701019b51ca754b6c782114825798973d8187d6",
		    "blockindex": 2,
		    "blocktime": 1533092896,
		    "txid": "11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		    "walletconflicts": [],
		    "time": 1533092879,
		    "timereceived": 1533092879,
		    "bip125-replaceable": "no",
		    "details": [
		      {
		        "account": "",
		        "category": "send",
		        "amount": 0,
		        "vout": 0,
		        "fee": -0.2012,
		        "abandoned": false
		      }
		    ],
		    "hex": "020000000159c0514feea50f915854d9ec45bc6458bb14419c78b17e7be3f7fd5f563475b5010000006a473044022072d64a1f4ea2d54b7b05050fc853ab192c91cc5ca17e23007867f92f2ab59d9202202b8c9ab9348c8edbb3b98b1788382c8f37642ec9bd6a4429817ab79927319200012103520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140feffffff02000000000000000063010403400d0301644440c10f190000000000000000000000006b22910b1e302cf74803ffd1691c2ecb858d3712000000000000000000000000000000000000000000000000000000000000000a14be528c8378ff082e4ba43cb1baa363dbf3f577bfc260e66272970100001976a9146b22910b1e302cf74803ffd1691c2ecb858d371288acb00f0000"
		  }
	*/
	GetTransactionResponse struct {
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
)

func (r *GetTransactionRequest) MarshalJSON() ([]byte, error) {
	/*
		1. "txid"                  (string, required) The transaction id
		2. "include_watchonly"     (bool, optional, default=false) Whether to include watch-only addresses in balance calculation and details[]
		3. "waitconf"              (int, optional, default=0) Wait for enough confirmations before returning
	*/
	return json.Marshal([]interface{}{
		r.Txid,
	})
}

func (r *GetTransactionResponse) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		return EmptyResponseErr
	}
	type Response GetTransactionResponse
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	*r = GetTransactionResponse(resp)

	return nil
}

// ========== SearchLogs ============= //

type (
	SearchLogsRequest struct {
		FromBlock *big.Int
		ToBlock   *big.Int
		Addresses []string
	}

	SearchLogsResponse []TransactionReceiptStruct
)

func (r *SearchLogsRequest) MarshalJSON() ([]byte, error) {
	/*
		1. "fromBlock"        (numeric, required) The number of the earliest block (latest may be given to mean the most recent block).
		2. "toBlock"          (string, required) The number of the latest block (-1 may be given to mean the most recent block).
		3. "address"          (string, optional) An address or a list of addresses to only get logs from particular account(s).
		4. "topics"           (string, optional) An array of values from which at least one must appear in the log entries. The order is important, if you want to leave topics out use null, e.g. ["null", "0x00..."].
		5. "minconf"          (uint, optional, default=0) Minimal number of confirmations before a log is returned
	*/
	return json.Marshal([]interface{}{
		r.FromBlock,
		r.ToBlock,
		map[string][]string{
			"addresses": r.Addresses,
		},
	})
}
