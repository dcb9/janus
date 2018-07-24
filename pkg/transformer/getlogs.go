package transformer

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/dcb9/go-ethereum/common/hexutil"
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
)

func (m *Manager) GetLogs(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	var params []json.RawMessage
	if err := unmarshalRequest(req.Params, &params); err != nil {
		return nil, err
	}

	var filter eth.GetLogsFilter
	if len(params) > 0 {
		if err := unmarshalRequest(params[0], &filter); err != nil {
			return nil, err
		}
	}

	if len(filter.Topics) != 0 {
		return nil, errors.New("topics is not supported yet")
	}

	from, err := m.getQtumBlockNumber(filter.FromBlock, 0)
	if err != nil {
		return nil, err
	}

	to, err := m.getQtumBlockNumber(filter.ToBlock, -1)
	if err != nil {
		return nil, err
	}

	var addresses []string
	if filter.Address != nil {
		if filter.Address[0] == '"' {
			var addr string
			if err = json.Unmarshal(filter.Address, &addr); err != nil {
				return nil, err
			}
			addresses = append(addresses, addr)
		} else {
			if err = json.Unmarshal(filter.Address, &addresses); err != nil {
				return nil, err
			}
		}
		for i, _ := range addresses {
			addresses[i] = RemoveHexPrefix(addresses[i])
		}
	}

	newParams, err := json.Marshal([]interface{}{
		from,
		to,
		map[string][]string{
			"addresses": addresses,
		},
	})

	req.Params = newParams
	req.Method = qtum.MethodSearchlogs

	return m.GetLogsResp, nil
}

func (m *Manager) GetLogsResp(result json.RawMessage) (interface{}, error) {
	var receipts []*qtum.TransactionReceipt
	if err := json.Unmarshal(result, &receipts); err != nil {
		return nil, err
	}

	logs := make([]eth.Log, 0)
	for _, receipt := range receipts {
		logs = append(logs, getEthLogs(receipt)...)
	}

	return logs, nil
}

func (m *Manager) getQtumBlockNumber(ethBlock json.RawMessage, defaultVal int64) (*big.Int, error) {
	if ethBlock == nil {
		return big.NewInt(defaultVal), nil
	}

	if isString(ethBlock) {
		var ethBlockStr string
		if err := json.Unmarshal(ethBlock, &ethBlockStr); err != nil {
			return nil, err
		}

		switch ethBlockStr {
		case "latest":
			return big.NewInt(-1), nil
		case "pending", "earliest":
			return nil, errors.New(`tags, "pending" and "earliest", are unsupported`)
		default:
			return hexutil.DecodeBig(ethBlockStr)
		}
	}

	var b int64
	if err := json.Unmarshal(ethBlock, &b); err != nil {
		return nil, err
	}

	return big.NewInt(b), nil
}

func isString(v json.RawMessage) bool {
	return v[0] == '"'
}

func getEthLogs(receipt *qtum.TransactionReceipt) []eth.Log {
	logs := make([]eth.Log, 0, len(receipt.Log))
	for index, log := range receipt.Log {
		topics := make([]string, 0, len(log.Topics))
		for _, topic := range log.Topics {
			topics = append(topics, AddHexPrefix(topic))
		}
		logs = append(logs, eth.Log{
			TransactionHash:  AddHexPrefix(receipt.TransactionHash),
			TransactionIndex: hexutil.EncodeUint64(receipt.TransactionIndex),
			BlockHash:        AddHexPrefix(receipt.BlockHash),
			BlockNumber:      hexutil.EncodeUint64(receipt.BlockNumber),
			Data:             AddHexPrefix(log.Data),
			Address:          AddHexPrefix(log.Address),
			Topics:           topics,
			LogIndex:         hexutil.EncodeUint64(uint64(index)),
		})
	}
	return logs
}

// Eth
// eth_getLogs
//fromBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
//toBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
//address: DATA|Array, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.
//topics: Array of DATA, - (optional) Array of 32 Bytes DATA topics. Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
//blockhash: DATA, 32 Bytes - (optional, future) With the addition of EIP-234, blockHash will be a new filter option which restricts the logs returned to the single block with the 32-byte hash blockHash. Using blockHash is equivalent to fromBlock = toBlock = the block number with hash blockHash. If blockHash is present in in the filter criteria, then neither fromBlock nor toBlock are allowed.

// Qtum
//+ qcli help searchlogs
//searchlogs <fromBlock> <toBlock> (address) (topics)
//requires -logevents to be enabled
//Argument:
//1. "fromBlock"        (numeric, required) The number of the earliest block (latest may be given to mean the most recent block).
//2. "toBlock"          (string, required) The number of the latest block (-1 may be given to mean the most recent block).
//3. "address"          (string, optional) An address or a list of addresses to only get logs from particular account(s).
//4. "topics"           (string, optional) An array of values from which at least one must appear in the log entries. The order is important, if you want to leave topics out use null, e.g. ["null", "0x00..."].
//5. "minconf"          (uint, optional, default=0) Minimal number of confirmations before a log is returned
//
//Examples:
//> qtum-cli searchlogs 0 100 '{"addresses": ["12ae42729af478ca92c8c66773a3e32115717be4"]}' '{"topics": ["null","b436c2bf863ccd7b8f63171201efd4792066b4ce8e543dde9c3e9e9ab98e216c"]}'
//> curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "searchlogs", "params": [0 100 {"addresses": ["12ae42729af478ca92c8c66773a3e32115717be4"]} {"topics": ["null","b436c2bf863ccd7b8f63171201efd4792066b4ce8e543dde9c3e9e9ab98e216c"]}] }' -H 'content-type: text/plain;' http://127.0.0.1:3889/
