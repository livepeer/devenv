require("babel-register")
require("babel-polyfill")

module.exports = {
    networks: {
        development: {
            host: "localhost",
            port: 8545,
            network_id: "*" // Match any network id
        },
        lpTestNet: {
            from: "0x94107cb2261e722f9f4908115546eeee17decada",
            host: "localhost",
            port: 8545,
            network_id: 54321,
            gas: 6700000
        }
    }
};
