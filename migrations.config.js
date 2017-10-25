module.exports = {
    bondingManager: {
        numActiveTranscoders: 5,
        unbondingPeriod: 2
    },
    jobsManager: {
        verificationRate: 10,
        jobEndingPeriod: 50,
        verificationPeriod: 50,
        slashingPeriod: 50,
        failedVerificationSlashAmount: 20,
        missedVerificationSlashAmount: 30,
        finderFee: 4
    },
    roundsManager: {
        blockTime: 1,
        roundLength: 5
    },
    faucet: {
        faucetAmount: 100000000000000000000,
        requestAmount: 1000000,
        requestWait: 2,
        whitelist: []
    },
    minter: {
        initialTokenSupply: 10000000 * Math.pow(10, 18),
        yearlyInflation: 26
    },
    verifier: {
        verificationCodeHash: "QmWdbVR8SUS9TU5a9HFP2qG18ck6Vh4mL2PYsHcB9sHXN7",
        solvers: ["0x0ddb225031ccb58ff42866f82d907f7766899014"],
        gasPrice: 20000000000,
        gasLimit: 3000000
    }
}
