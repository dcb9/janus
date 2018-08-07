module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
  networks: {
    development: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*",
      from: "0xd617b07d7ede55244246a807d22aa5e705e13301"
    },
    qtum: {
      host: "127.0.0.1",
      port: 33889,
      network_id: "*",
      // acc1 =
      from: "0x7926223070547d2d15b2ef5e7383e541c338ffe9",
      // acc2 = "0x3a895d2af552600f1f585425318c13a5aa25f01a"

      // gas: ,
      gasPrice: "0x64"
    }
  }
};
