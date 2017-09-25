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
        roundLength: 5,
    },
    faucet: {
        faucetAmount: 100000000000000000000,
        requestAmount: 1000000,
        requestWait: 2,
        whitelist: []
    }
}
