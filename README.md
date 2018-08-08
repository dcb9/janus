Qtum adapter to Ethereum JSON RPC
=====================================

Table of Contents
=================
* [Start server](#start-server)
* [ERC20 with QtumJS](#erc20-with-qtumjs)
   * [Deploy myToken](#deploy-mytoken)
   * [Methods](#methods)
      * [mint](#mint)
      * [balance](#balance)
      * [transfer](#transfer)
      * [logs](#logs)
      * [events](#events)
* [Interact with QtumJS](#interact-with-qtumjs)
* [ERC20 With QtumJS](#erc20-with-qtumjs-1)
* [Try to interact with contract](#try-to-interact-with-contract)
   * [Assumption parameters](#assumption-parameters)
   * [createcontract method](#createcontract-method)
   * [gettransaction method](#gettransaction-method)
   * [gettransactionreceipt method](#gettransactionreceipt-method)
   * [sendtocontract method](#sendtocontract-method)
   * [callcontract method](#callcontract-method)
   * [sendtoaddress method](#sendtoaddress-method)
* [Support ETH methods](#support-eth-methods)
* [Todo list](#todo-list)
* [Known issues](#known-issues)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

## Start server

```
$ go run cli/janus/main.go --qtum-rpc="http://qtum:test@localhost:3889" --port=23889  --dev
```

Open another terminal to do some fancy stuff

// prepare two accounts

```
/dapp # qcli getnewaddress
qc614qvV6Jgh6yei3wtCyWneqQtsq8Cpng

/dapp # qcli getnewaddress
qd73nuEB1wFsz24TcQeg1QQU9zeQB7Uu1h

/dapp # qcli gethexaddress qc614qvV6Jgh6yei3wtCyWneqQtsq8Cpng
cb3cb8375fe457a11f041f9ff55373e1a5a78d19

/dapp # qcli gethexaddress qd73nuEB1wFsz24TcQeg1QQU9zeQB7Uu1h
d66789418ca152f5720b1c8dd04e9ff2f3891f6f
```

So your hexed addresses are `0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19` and `0xd66789418ca152f5720b1c8dd04e9ff2f3891f6f`

```
// set environment

$ cd playground
$ yarn install
$ export ETH_RPC=http://0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19:@localhost:23889
```

## ERC20 with QtumJS

### Deploy myToken
```
$ sh deploy-myToken.sh
  + solar deploy openzeppelin-solidity/contracts/token/ERC20/CappedToken.sol --gasPrice=0.0000001 '[21000000]' --force
  exec: solc [openzeppelin-solidity/contracts/token/ERC20/CappedToken.sol --combined-json bin,metadata --optimize --allow-paths /Users/bob/Documents/golangWorkspace/src/github.com/dcb9/janus/playground]
  cli gasPrice 0.0000001 1e-07
  gasPrice 1e-07 100
  gasPriceWei 100
  txHash: 0xfe30795aafce57532af6215fe57ed3382f43fbeaf2b3c42c8ea4d9f34ab0be55
  contractAddress: 0x90f3e8062c8537ee4825fd384caef0260795f8df
  🚀  All contracts confirmed
     deployed openzeppelin-solidity/contracts/token/ERC20/CappedToken.sol => 0x90f3e8062c8537ee4825fd384caef0260795f8df
```

### Methods

#### mint

```
$ node myToken.js mint 0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19 100

mint tx: undefined
{ hash: '0x5ee0463596cf35c00363f10b5781499aa6693e0477965053b5537716e84113c6',
  nonce: '',
  blockHash: '0x',
  blockNumber: '',
  transactionIndex: '',
  from: '',
  to: '',
  value: '0x0',
  gasPrice: '0x64',
  gas: '0x30d40',
  input: '0x40c10f19000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d190000000000000000000000000000000000000000000000000000000000000064',
  method: 'mint',
  confirm: [Function: confirm] }
✔ confirm mint
tx receipt: {
  "transactionHash": "0x5ee0463596cf35c00363f10b5781499aa6693e0477965053b5537716e84113c6",
  "transactionIndex": "0x2",
  "blockHash": "0xe7e8523ff95cd2f6663992e2c4e80d354fb20d1c67d6ffbbbb0c1448758f61a1",
  "blockNumber": "0x353f",
  "cumulativeGasUsed": "0x10bdc",
  "gasUsed": "0x10bdc",
  "contractAddress": "0x90f3e8062c8537ee4825fd384caef0260795f8df",
  "logsBloom": "",
  "status": "0x1",
  "from": "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19",
  "to": "0x90f3e8062c8537ee4825fd384caef0260795f8df",
  "logs": [
    {
      "amount": "64",
      "to": "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19",
      "type": "Mint"
    },
    {
      "value": "64",
      "from": "0x0000000000000000000000000000000000000000",
      "to": "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19",
      "type": "Transfer"
    }
  ],
  "rawlogs": [
    {
      "logIndex": "0x0",
      "transactionIndex": "0x2",
      "transactionHash": "0x5ee0463596cf35c00363f10b5781499aa6693e0477965053b5537716e84113c6",
      "blockHash": "0xe7e8523ff95cd2f6663992e2c4e80d354fb20d1c67d6ffbbbb0c1448758f61a1",
      "blockNumber": "0x353f",
      "address": "0x90f3e8062c8537ee4825fd384caef0260795f8df",
      "data": "0x0000000000000000000000000000000000000000000000000000000000000064",
      "topics": [
        "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885",
        "0x000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d19"
      ]
    },
    {
      "logIndex": "0x1",
      "transactionIndex": "0x2",
      "transactionHash": "0x5ee0463596cf35c00363f10b5781499aa6693e0477965053b5537716e84113c6",
      "blockHash": "0xe7e8523ff95cd2f6663992e2c4e80d354fb20d1c67d6ffbbbb0c1448758f61a1",
      "blockNumber": "0x353f",
      "address": "0x90f3e8062c8537ee4825fd384caef0260795f8df",
      "data": "0x0000000000000000000000000000000000000000000000000000000000000064",
      "topics": [
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "0x0000000000000000000000000000000000000000000000000000000000000000",
        "0x000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d19"
      ]
    }
  ]
}
```

#### balance

```
$  node myToken.js balance 0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19

balance: 100
```

#### transfer

```
$ node myToken.js transfer 0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19  0xd66789418ca152f5720b1c8dd04e9ff2f3891f6f 5

transfer tx: undefined
{ hash: '0x0cfe192b3244bc0eb1089a72962d2f164156ff40bea46e160c74fc86b0403d9a',
  nonce: '',
  blockHash: '0x',
  blockNumber: '',
  transactionIndex: '',
  from: '',
  to: '',
  value: '0x0',
  gasPrice: '0x64',
  gas: '0x30d40',
  input: '0xa9059cbb000000000000000000000000d66789418ca152f5720b1c8dd04e9ff2f3891f6f0000000000000000000000000000000000000000000000000000000000000005',
  method: 'transfer',
  confirm: [Function: confirm] }
✔ confirm transfer
```

#### logs

```
$ node myToken.js logs

[
  {
    "logIndex": "0x0",
    "transactionIndex": "0x2",
    "transactionHash": "0x5ee0463596cf35c00363f10b5781499aa6693e0477965053b5537716e84113c6",
    "blockHash": "0xe7e8523ff95cd2f6663992e2c4e80d354fb20d1c67d6ffbbbb0c1448758f61a1",
    "blockNumber": "0x353f",
    "address": "0x90f3e8062c8537ee4825fd384caef0260795f8df",
    "data": "0x0000000000000000000000000000000000000000000000000000000000000064",
    "topics": [
      "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885",
      "0x000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d19"
    ],
    "event": {
      "amount": "64",
      "to": "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19",
      "type": "Mint"
    }
  },
  {
    "logIndex": "0x1",
    "transactionIndex": "0x2",
    "transactionHash": "0x5ee0463596cf35c00363f10b5781499aa6693e0477965053b5537716e84113c6",
    "blockHash": "0xe7e8523ff95cd2f6663992e2c4e80d354fb20d1c67d6ffbbbb0c1448758f61a1",
    "blockNumber": "0x353f",
    "address": "0x90f3e8062c8537ee4825fd384caef0260795f8df",
    "data": "0x0000000000000000000000000000000000000000000000000000000000000064",
    "topics": [
      "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
      "0x0000000000000000000000000000000000000000000000000000000000000000",
      "0x000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d19"
    ],
    "event": {
      "value": "64",
      "from": "0x0000000000000000000000000000000000000000",
      "to": "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19",
      "type": "Transfer"
    }
  },
  {
    "logIndex": "0x0",
    "transactionIndex": "0x2",
    "transactionHash": "0x0cfe192b3244bc0eb1089a72962d2f164156ff40bea46e160c74fc86b0403d9a",
    "blockHash": "0x307c14f3a536b62e648566bec19f86270098a3d97d35f0d0cbaed4500150cb80",
    "blockNumber": "0x3544",
    "address": "0x90f3e8062c8537ee4825fd384caef0260795f8df",
    "data": "0x0000000000000000000000000000000000000000000000000000000000000005",
    "topics": [
      "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
      "0x000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d19",
      "0x000000000000000000000000d66789418ca152f5720b1c8dd04e9ff2f3891f6f"
    ],
    "event": {
      "value": "5",
      "from": "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19",
      "to": "0xd66789418ca152f5720b1c8dd04e9ff2f3891f6f",
      "type": "Transfer"
    }
  }
]
```

#### events

```
$ node myToken.js events

Subscribed to contract events
Ctrl-C to terminate events subscription
{ logIndex: '0x0',
  transactionIndex: '0x2',
  transactionHash: '0xe1e8afd1591bb4ef110fe4ddddf7de2bc1c04bbace3eb079cb95c8f8c5214729',
  blockHash: '0xd4e8fcea409a82c303823faee8164a7e5e57531c8cf50d37082fd3c128fb1e62',
  blockNumber: '0x3549',
  address: '0x90f3e8062c8537ee4825fd384caef0260795f8df',
  data: '0x0000000000000000000000000000000000000000000000000000000000000005',
  topics:
   [ '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef',
     '0x000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d19',
     '0x000000000000000000000000d66789418ca152f5720b1c8dd04e9ff2f3891f6f' ],
  event:
   Result {
     value: <BN: 5>,
     from: '0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19',
     to: '0xd66789418ca152f5720b1c8dd04e9ff2f3891f6f',
     type: 'Transfer' } }
```


## Interact with QtumJS

```
$ sh deploy-SimpleStore.sh
  + solar deploy ./contracts/SimpleStore.sol '["1"]' --gasPrice=0.0000001 --force
  exec: solc [./contracts/SimpleStore.sol --combined-json bin,metadata --optimize --allow-paths /Users/bob/Documents/golangWorkspace/src/github.com/dcb9/janus/playground]
  cli gasPrice 0.0000001 1e-07
  gasPrice 1e-07 100
  gasPriceWei 100
  txHash: 0x95472d05243864764211bd8c6d8110fa397bd045cff78d845c1250bdff789bc7
  contractAddress: 0x6997a4803d75964b8d093a939c227a16833d23ad
  🚀  All contracts confirmed
     deployed ./contracts/SimpleStore.sol => 0x6997a4803d75964b8d093a939c227a16833d23ad

$ node test-SimpleStore.js
exec: await simpleStoreContract.call("get", [], {gasPrice: 100})
call { rawResult: '0x0000000000000000000000000000000000000000000000000000000000000001',
  outputs: [ <BN: 1> ],
  logs: [] }

exec: await simpleStoreContract.send("set", [82009999], {gasPrice: 100})
tx { hash: '0x23a0d715ef4fc2ce8bcf79bf1427e3fea6af38905efab9668672e693591f3ee4',
  nonce: '',
  blockHash: '0x',
  blockNumber: '',
  transactionIndex: '',
  from: '',
  to: '',
  value: '0x0',
  gasPrice: '0x64',
  gas: '0x30d40',
  input: '0x60fe47b10000000000000000000000000000000000000000000000000000000004e35f8f',
  method: 'set',
  confirm: [Function: confirm] }

exec: await tx.confirm(0)
receipt { transactionHash: '0x23a0d715ef4fc2ce8bcf79bf1427e3fea6af38905efab9668672e693591f3ee4',
  transactionIndex: '0x2',
  blockHash: '0x6b8273375b3a8dff6701c4151d03aa2e3211fbb3f2bea558d16a762fe0cd2b1a',
  blockNumber: '0x2eaf',
  cumulativeGasUsed: '0x702e',
  gasUsed: '0x702e',
  contractAddress: '0x6997a4803d75964b8d093a939c227a16833d23ad',
  logsBloom: '',
  status: '0x1',
  from: '0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19',
  to: '0x6997a4803d75964b8d093a939c227a16833d23ad',
  logs:
   [ Result {
       from: '0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19',
       _oldValue: <BN: 1>,
       _newValue: <BN: 4e35f8f>,
       type: 'UpdateValue' } ],
  rawlogs:
   [ { logIndex: '0x0',
       transactionIndex: '0x2',
       transactionHash: '0x23a0d715ef4fc2ce8bcf79bf1427e3fea6af38905efab9668672e693591f3ee4',
       blockHash: '0x6b8273375b3a8dff6701c4151d03aa2e3211fbb3f2bea558d16a762fe0cd2b1a',
       blockNumber: '0x2eaf',
       address: '0x6997a4803d75964b8d093a939c227a16833d23ad',
       data: '0x000000000000000000000000cb3cb8375fe457a11f041f9ff55373e1a5a78d190000000000000000000000000000000000000000000000000000000000000001',
       topics: [Array] } ] }

exec: await simpleStoreContract.call("get", [], {gasPrice: 100})
call { rawResult: '0x0000000000000000000000000000000000000000000000000000000004e35f8f',
  outputs: [ <BN: 4e35f8f> ],
  logs: [] }

```

## ERC20 With QtumJS

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
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19","gas":"0x6691b7","gasPrice":"0x174876e800","value":"0xffffff", "to": "0xd66789418ca152f5720b1c8dd04e9ff2f3891f6f"}]}' \
     'http://localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0x978ed14c122dca1669df875e2cc33302a6edd13b7a8a5a30e3a53ef53b53bbf4",
  "id": "10"
}


$ curl --header 'Content-Type: application/json' --data \
       '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x978ed14c122dca1669df875e2cc33302a6edd13b7a8a5a30e3a53ef53b53bbf4"]}' \
  'localhost:23889'

// notice: the tx receipt of sendtoaddress is an empty array
{
  "jsonrpc": "2.0",
  "result": [],
  "id": "10"
}

$ curl --header 'Content-Type: application/json' --data \
       '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x978ed14c122dca1669df875e2cc33302a6edd13b7a8a5a30e3a53ef53b53bbf4"]}' \
  'localhost:23889'

// notice: blockNumber, transactionIndex, from, to are empty, because tx receipt of sendtoaddress is an empty array
{
  "jsonrpc": "2.0",
  "result": {
    "hash": "0x978ed14c122dca1669df875e2cc33302a6edd13b7a8a5a30e3a53ef53b53bbf4",
    "nonce": "",
    "blockHash": "0x9a5c002cac26df0bbd77099412dff3bd542741a1bb6e955cc161b76a83b8626f",
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
curl --header 'Content-Type: application/json' --data \
       '{"id":"10","jsonrpc":"2.0","method":"eth_getLogs","params":[]}' \
  'localhost:23889'


## Support ETH methods

- eth_sendTransaction
- eth_call
- eth_getTransactionByHash
- eth_getTransactionReceipt
- eth_blockNumber
- net_version
  - returns string // current network name as defined in BIP70 (main, test, regtest)
- eth_getLogs
  - topics is not supported yet
  - tags, "pending" and "earliest", are unsupported
- eth_accounts
- eth_getCode
- eth_newBlockFilter
- eth_getFilterChanges
  - only support filters created with `eth_newBlockFilter`
- eth_uninstallFilter

## Known issues

- eth_getTransactionReceipt
  - `logsBloom` is an empty string
  - result will be an empty array if the txid of the transaction is a transfer operation
- eth_getTransactionByHash
  - `nonce` is an empty string
  - `blockNumber`, `transactionIndex`, `from`, `to`, `value` will be empty, if the txid of the transaction is a transfer operation
- eth_accounts
  - only return addresses which are linked to default account
