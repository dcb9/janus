Qtum adapter to Ethereum JSON RPC
=====================================

## Start server

```
$ go run cli/janus/main.go --qtum-rpc="http://qtum:test@localhost:3889" --port=23889  --dev
```

## Try to interact with contract

see: [Qtum smart contract](http://book.qtum.site/en/part4/smart-contract.html)

### Assumption parameters

Assumed that you have a **contract** like this:

```solidity
pragma solidity ^0.4.18;

contract SimpleStore {
  constructor(uint _value) public {
    value = _value;
  }

  function set(uint newValue) public {
    value = newValue;
  }

  function get() public constant returns (uint) {
    return value;
  }

  uint value;
}
```

so that the **bytecode** is

```
solc --optimize --bin contracts/SimpleStore.sol

======= contracts/SimpleStore.sol:SimpleStore =======
Binary:
608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b810029
```

**constructor parameters** is `0000000000000000000000000000000000000000000000000000000000000001`
**account address**
```
/dapp # qcli gethexaddress qc614qvV6Jgh6yei3wtCyWneqQtsq8Cpng
cb3cb8375fe457a11f041f9ff55373e1a5a78d19
```

### createcontract method

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19","gas":"0x6691b7","gasPrice":"0x174876e800","data":"0x608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b8100290000000000000000000000000000000000000000000000000000000000000001"}]}' \
     'http://localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0x6da39dc909debf70a536bbc108e2218fd7bce23305ddc00284075df5dfccc21b",
  "id": "10"
}
```

### gettransaction method

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x6da39dc909debf70a536bbc108e2218fd7bce23305ddc00284075df5dfccc21b"]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": {
    "hash": "0x6da39dc909debf70a536bbc108e2218fd7bce23305ddc00284075df5dfccc21b",
    "nonce": "",
    "blockHash": "0xa5f0db33370d6a3e83ace9ed2b3ff74c29ad70b78427eb67de1d959dfa485085",
    "blockNumber": "0x1c51",
    "transactionIndex": "0x2",
    "from": "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19",
    "to": "0x1d96667c8de1a6d8a2a393d6518f376ed3239dd3",
    "value": "0x0",
    "gasPrice": "0x28",
    "gas": "0x6691b7",
    "input": "0x608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b8100290000000000000000000000000000000000000000000000000000000000000001"
  },
  "id": "10"
}
```

### gettransactionreceipt method

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x6da39dc909debf70a536bbc108e2218fd7bce23305ddc00284075df5dfccc21b"]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": {
    "transactionHash": "0x6da39dc909debf70a536bbc108e2218fd7bce23305ddc00284075df5dfccc21b",
    "transactionIndex": "0x2",
    "blockHash": "0xa5f0db33370d6a3e83ace9ed2b3ff74c29ad70b78427eb67de1d959dfa485085",
    "blockNumber": "0x1c51",
    "cumulativeGasUsed": "0x1e8a9",
    "gasUsed": "0x1e8a9",
    "contractAddress": "0x1d96667c8de1a6d8a2a393d6518f376ed3239dd3",
    "logs": [],
    "logsBloom": "",
    "status": "0x1"
  },
  "id": "10"
}
```

### sendtocontract method

the ABI code of set method with param '["2"]' is `60fe47b10000000000000000000000000000000000000000000000000000000000000002`

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19","gas":"0x6691b7","gasPrice":"0x174876e800","to":"0x1d96667c8de1a6d8a2a393d6518f376ed3239dd3","data":"60fe47b10000000000000000000000000000000000000000000000000000000000000002"}]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0xb6a315733207992115e8aa002b7b9543d34839f7265f3f5399453ebf54febe71",
  "id": "10"
}
```

### callcontract method

get method's ABI code is `6d4ce63c`

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_call","params":[{"from":"0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19","gas":"0x6691b7","gasPrice":"0x174876e800","to":"0x1d96667c8de1a6d8a2a393d6518f376ed3239dd3","data":"6d4ce63c"},"latest"]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0x0000000000000000000000000000000000000000000000000000000000000002",
  "id": "10"
}
```

### sendtoaddress method

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0xcf4b32865abe674c9e75c568bc30cf3cbc701720","gas":"0x6691b7","gasPrice":"0x174876e800","value":"0xfffffff", "to": "0xc6f125a5a7ac965f22ef220d8dd639fe44c68db7"}]}' \
     'http://localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0xf28eb8ff36735920db213b2db6693c6baa2a5642b4e7860b88f7dc1204ac4ca3",
  "id": "10"
}


$ curl --header 'Content-Type: application/json' --data \
       '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xf28eb8ff36735920db213b2db6693c6baa2a5642b4e7860b88f7dc1204ac4ca3"]}' \
  'localhost:23889'

// notice: the tx receipt of sendtoaddress is an empty array
{
  "jsonrpc": "2.0",
  "result": [],
  "id": "10"
}


$ curl --header 'Content-Type: application/json' --data \
       '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xf28eb8ff36735920db213b2db6693c6baa2a5642b4e7860b88f7dc1204ac4ca3"]}' \
  'localhost:23889'

// notice: blockNumber, transactionIndex, from, to are empty, because tx receipt of sendtoaddress is an empty array
{
  "jsonrpc": "2.0",
  "result": {
    "hash": "0xf28eb8ff36735920db213b2db6693c6baa2a5642b4e7860b88f7dc1204ac4ca3",
    "nonce": "",
    "blockHash": "0x131a22facf5e1bb4ce0618ac8a2e4d4ee8a1a3602e45ffacc577c941762174f4",
    "blockNumber": "",
    "transactionIndex": "",
    "from": "",
    "to": "",
    "value": "0x0",
    "gasPrice": "",
    "gas": "",
    "input": ""
  },
  "id": "10"
}
```

## Support ETH methods

- eth_sendTransaction
- eth_call
- eth_getTransactionByHash
- eth_getTransactionReceipt

## Todo list

## Known issues

- [ ] eth_getTransactionReceipt
  - [ ] `logsBloom` is an empty string
  - [ ] result will be an empty array if the txid of the transaction is a transfer operation
- [ ] eth_getTransactionByHash
  - [ ] `nonce` is an empty string
  - [ ] `blockNumber`, `transactionIndex`, `from`, `to`, `value` will be empty, if the txid of the transaction is a transfer operation

