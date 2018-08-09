const url = require('url');

const rpcURL =  process.env.ETH_RPC;
const qtumAccount  = url.parse(rpcURL).auth.split(":")[0]

// assume: node 8 or above
const ora = require("ora")
const parseArgs = require("minimist")

const qtum = require("qtumjs")
const rpc = new qtum.EthRPC(rpcURL, qtumAccount)
const repoData = require("./solar.development.json")
const {
  sender,
  ...info
} = repoData.contracts['contracts/MyToken.sol']
const myToken = new qtum.Contract(rpc, info)

const opts = {gasPrice: 100}

async function totalSupply() {
  const result = await myToken.call("totalSupply", [], opts)

  // supply is a BigNumber instance (see: bn.js)
  const supply = result.outputs[0]

  console.log("supply", supply.toNumber())
}

async function balanceOf(owner) {
  const res = await myToken.call("balanceOf", [owner], opts)

  // balance is a BigNumber instance (see: bn.js)
  const balance = res.outputs[0]

  console.log(`balance:`, balance.toNumber())
}

async function mint(toAddr, amount) {
  const tx = await myToken.send("mint", [toAddr, amount], opts)

  console.log("mint tx:", tx.txid)
  console.log(tx)

  // or: await tx.confirm(1)
  const confirmation = tx.confirm(1)
  ora.promise(confirmation, "confirm mint")
  const receipt = await confirmation
  console.log("tx receipt:", JSON.stringify(receipt, null, 2))
}

async function transfer(fromAddr, toAddr, amount) {
  const tx = await myToken.send("transfer", [toAddr, amount], {
    senderAddress: fromAddr,
    gasPrice: 100,
  })

  console.log("transfer tx:", tx.txid)
  console.log(tx)

  // or: await tx.confirm(1)
  const confirmation = tx.confirm(1)
  ora.promise(confirmation, "confirm transfer")
  await confirmation
}

async function streamEvents() {
  console.log("Subscribed to contract events")
  console.log("Ctrl-C to terminate events subscription")

  myToken.onLog((entry) => {
    console.log(entry)
  }, { minconf: 1 })
}

async function getLogs(fromBlock, toBlock) {
  const logs = await myToken.logs({
    fromBlock,
    toBlock,
    minconf: 1,
  })

  console.log(JSON.stringify(logs, null, 2))
}

async function main() {
  const argv = parseArgs(process.argv.slice(2), {"string": '_'})

  const cmd = argv._[0]

  if (process.env.DEBUG) {
    console.log("argv", argv)
    console.log("cmd", cmd)
  }

  switch (cmd) {
    case "supply":
    case "totalSupply":
      await totalSupply()
      break
    case "balance":
      const owner = argv._[1]
      if (!owner) {
        throw new Error("please specify an address")
      }
      await balanceOf(owner)
      break
    case "mint":
      const mintToAddr = argv._[1]
      const mintAmount = parseInt(argv._[2])
      await mint(mintToAddr, mintAmount)

      break
    case "transfer":
      const fromAddr = argv._[1]
      const toAddr = argv._[2]
      const amount = argv._[3]

      await transfer(fromAddr, toAddr, amount)
      break
    case "logs":
      const fromBlock = parseInt(argv._[1]) || 0
      const toBlock = parseInt(argv._[2]) || "latest"

      await getLogs(fromBlock, toBlock)
      break
    case "events":
      await streamEvents() // logEvents will never return
      break
    default:
      console.log("unrecognized command", cmd)
  }
}

main().catch(err => {
  console.log("error", err)
})
