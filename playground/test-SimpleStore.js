// YOUR_QTUM_ACCOUNT
const qtumAccount = "0xcb3cb8375fe457a11f041f9ff55373e1a5a78d19"

const qtum = require("qtumjs")
const rpc = new qtum.EthRPC(`http://${qtumAccount}:@localhost:23889`, qtumAccount)
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
