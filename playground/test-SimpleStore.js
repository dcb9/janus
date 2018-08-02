// YOUR_QTUM_ACCOUNT
const url = require('url');

const rpcURL=  process.env.ETH_RPC;
const qtumAccount  = url.parse(rpcURL).auth.split(":")[0]

const qtum = require("qtumjs")
const rpc = new qtum.EthRPC(rpcURL, qtumAccount)
const repoData = require("./solar.development.json")
const {
  sender,
  ...info
} = repoData.contracts['./contracts/SimpleStore.sol']
const simpleStoreContract = new qtum.Contract(rpc, info)

const opts = {gasPrice: 100}


async function test() {
  console.log('exec: await simpleStoreContract.call("get", [], {gasPrice: 100})')
  console.log("call", await simpleStoreContract.call("get", [], opts))
  console.log()

  const newVal = Math.floor((Math.random() * 100000000) + 1);
  console.log(`exec: await simpleStoreContract.send("set", [${newVal}], {gasPrice: 100})`)
  const tx = await simpleStoreContract.send("set", [newVal], opts)
  console.log("tx", tx)
  console.log()

  console.log('exec: await tx.confirm(0)')
  const receipt = await tx.confirm(0)
  console.log("receipt", receipt)
  console.log()

  console.log('exec: await simpleStoreContract.call("get", [], {gasPrice: 100})')
  console.log("call", await simpleStoreContract.call("get", [], opts))
  console.log()
}

test()
