#!/bin/bash

##
# The app formerly known as "testenv-cli"
# Purpose: Set up local Livepeer dev and testing environment.

##
srcDir=${LPSRC:-$HOME/src}
binDir=${LPBIN:-livepeer}
nodeBaseDataDir=$HOME/.lpdata

gethDir=${LPETH:-$HOME/.ethereum}
gethIPC=$gethDir/geth.ipc
gethPid=0
gethRunning=false
gethMiningAccount=
accountAddress=

protocolBuilt=false
protocolBranch="repo not found"
controllerAddress=

broadcasterRtmpPort=1935
broadcasterApiPort=8935
broadcasterPid=0
broadcasterRunning=false
broadcasterGeth=

transcoderRtmpPort=1936
transcoderApiPort=8936
transcoderPid=0
transcoderRunning=false
transcoderGeth=

##
#
# TODO: create separate commands
# $ lpdev geth [ init run reset ]
# $ lpdev protocol [ init deploy reset ]
# $ lpdev node [ broadcaster transcoder reset ]
# $ lpdev [ status wizard reset ]
#
##

##
# Display the status of the current environment
##

function __lpdev_status {
  echo "== Current Status ==
  "

  ##
  # Is geth set up and running?
  ##
  __lpdev_geth_refresh_status

  echo "Geth miner is running: $gethRunning ($gethPid)"
  if [ $gethRunning ]
  then
    gethAccounts=($(geth account list | cut -d' ' -f3 | tr -cd '[:alnum:]\n'))
    echo "Geth accounts:"
    for i in ${!gethAccounts[@]}
    do
      accountAddress=${gethAccounts[$i]}
      if [ $i -eq 0 ]
      then
        accountAddress="$accountAddress (miner)"
      fi
      echo "  $accountAddress"
    done
  fi

  echo ""

  ##
  # Is the protocol compiled and deployed?
  ##
  __lpdev_protocol_refresh_status

  echo "Protocol is built: $protocolBuilt (current branch: $protocolBranch)"

  if [ $controllerAddress ]
  then
    echo "Protocol deployed to: $controllerAddress"
  fi

  echo ""

  ##
  # Are nodes running?
  ##
  __lpdev_node_refresh_status

  echo "Broadcaster node is running: $broadcasterRunning ($broadcasterPid)"
  echo "Transcoder node is running: $transcoderRunning ($transcoderPid)"

  echo "
--
  "

}

function __lpdev_reset {

  echo "This will reset the dev environment"
  read -p "Are you sure you want to continue? [y/N] " -n 1 -r
  echo ""
  if [[ ! $REPLY =~ ^[Yy]$ ]]
  then
    return 1
  fi

  __lpdev_geth_reset
  __lpdev_protocol_reset
  __lpdev_node_reset
  echo "Local dev environment has been reset"

}

##
#
# Geth commands: init run reset
#
##

function __lpdev_geth_refresh_status {

  gethPid=$(pgrep -f "geth.*-mine")

  if [ -e $gethIPC ]
  then
    gethMiningAccount=$(geth attach ipc:/home/ubuntu/.ethereum/geth.ipc --exec "eth.coinbase" 2> /dev/null | grep "0x" | cut -d"x" -f2 | tr -cd "[:alnum:]")
  fi

  if [ -n "${gethPid}" ] && [ -n "${gethMiningAccount}" ]
  then
    gethRunning=true
  else
    gethRunning=false
  fi

}

function __lpdev_geth_reset {

  pkill -9 geth
  echo "Removing $gethDir and ~/.ethash"
  rm -rf $gethDir ~/.ethash
  unset gethPid
  unset gethMiningAccount
  gethRunning=false

}

function __lpdev_geth_init {

  __lpdev_geth_refresh_status

  if [ -n "${gethMiningAccount}" ]
  then
    echo "Geth mining account exists"
    return 1
  else
    echo "Creating miner account"
    gethMiningAccount=$(geth account new --password <(echo "") | cut -d' ' -f2 | tr -cd '[:alnum:]')
    echo "Created mining account $gethMiningAccount"
  fi

  if [ -d $gethDir/geth/chaindata ]
  then
    echo "Geth genesis was initialized"
    return 1
  fi

  echo "Setting up Geth data at ~/.ethereum"
  geth init <( cat << EOF
  {
    "config": {
      "chainId": 54321,
      "homesteadBlock": 1,
      "eip150Block": 2,
      "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "eip155Block": 3,
      "eip158Block": 3,
      "byzantiumBlock": 4,
      "clique": {
        "period": 2,
        "epoch": 30000
      }
    },
    "nonce": "0x0",
    "timestamp": "0x59bc2eff",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000${gethMiningAccount}0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "gasLimit": "0x663be0",
    "difficulty": "0x1",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "coinbase": "0x0000000000000000000000000000000000000000",
    "alloc": {
      "$gethMiningAccount": {
        "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
      }
    },
    "number": "0x0",
    "gasUsed": "0x0",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
  }
EOF
)

  if [ $? -ne 0 ]
  then
    echo "Could not initialize Geth"
  fi

}

function __lpdev_geth_run {

  __lpdev_geth_refresh_status

  if $gethRunning
  then
    echo "Geth is running, please kill it ($gethPid) or reset the environment if you'd like a fresh start."
    return 1
  fi

  echo "Running Geth miner with the following command:
  geth -networkid 54321
       -rpc
       -rpcapi 'personal,account,eth,web3,net'
       -targetgaslimit 6700000
       -unlock $gethMiningAccount
       --password <(echo \"\")
       -mine"

  nohup geth -networkid 54321 -rpc -rpcapi 'personal,account,eth,web3,net' -targetgaslimit 6700000 -unlock $gethMiningAccount --password <(echo "") -mine &>>$nodeBaseDataDir/geth.log &

  if [ $? -ne 0 ]
  then
    echo "Could not start Geth"
  else
    echo "Geth started successfully"
    disown
  fi

}

function __lpdev_protocol_refresh_status {

  if [ -d $srcDir/protocol/build ]
  then
    protocolBuilt=true
  fi

  if [ -d $srcDir/protocol ]
  then
    protocolBranch=$(cd $srcDir/protocol && git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/')
    controllerAddress=$(cd $srcDir/protocol && truffle networks | grep Controller | cut -d':' -f2 | tr -cd '[:alnum:]')
  fi

}

function __lpdev_protocol_reset {

  if [ -d $srcDir/protocol/build ]
  then
    rm -rf $srcDir/protocol/build
  fi

  protocolBuilt=false
  protocolBranch="repo not found"
  unset controllerAddress

}

function __lpdev_protocol_init {

  if [ -d $srcDir/protocol ]
  then
    echo "Protocol src directory exists"
  else
    echo "Cloning github.com/livepeer/protocol into src directory"
    OPWD=$PWD
    cd $srcDir
    git clone -b develop "https://github.com/livepeer/protocol.git"
    cd $OPWD
  fi

  ##
  # Install local dev truffle.js
  ##

  __lpdev_geth_refresh_status

  if [ -z "${gethMiningAccount}" ]
  then
    echo "Geth Mining Account not found"
    return 1
  fi

  if grep -q ${gethMiningAccount:-"none"} $srcDir/protocol/truffle.js
  then
    echo "Local dev version of $srcDir/protocol/truffle.js already exists"
  else
    echo "Installing local dev version of $srcDir/protocol/truffle.js"

    cat << EOF > $srcDir/protocol/truffle.js
module.exports = {
    networks: {
        development: {
            host: "localhost",
            port: 8545,
            network_id: "*" // Match any network id
        },
        lpTestNet: {
            from: "0x$gethMiningAccount",
            host: "localhost",
            port: 8545,
            network_id: 54321,
            gas: 6700000
        }
    }
};
EOF

    echo "Installing local dev version of $srcDir/protocol/migrations/migrations.config.js"

    cat << EOF > $srcDir/protocol/migrations/migrations.config.js
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
EOF

  fi

  ##
  # Update npm
  ##

  if [ -d $srcDir/protocol/node_modules ]
  then
    echo "Npm packages already installed"
  else
    echo "Running \`npm install\`"
    OPWD=$PWD
    cd $srcDir/protocol
    npm install
    cd $OPWD
  fi
}

function __lpdev_protocol_deploy {

  __lpdev_geth_refresh_status

  if ! $gethRunning
  then
    echo "Geth is not running, please start it before deploying protocol"
    return 1
  fi

  __lpdev_protocol_refresh_status

  if $protocolBuilt && [ -n "${controllerAddress}" ]
  then
    echo "Protocol already deployed ($controllerAddress)"
    read -p "Would you like to recompile and redeploy? [y/N] " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]
    then
      return 1
    fi
    redeployed=true
  fi

  OPWD=$PWD
  cd $srcDir/protocol
  echo "Running \`truffle migrate --reset --network lpTestNet\`"
  truffle migrate --reset --network lpTestNet
  cd $OPWD

  if $redeployed
  then
    echo "Don't forget to restart any nodes that may be using the previous controllerAddr!"
  fi
}

function __lpdev_node_refresh_status {

  broadcasterPid=$(pgrep -f "livepeer.*$broadcasterApiPort")
  if [ -n "${broadcasterPid}" ]
  then
    broadcasterRunning=true
  else
    broadcasterRunning=false
  fi

  transcoderPid=$(pgrep -f "livepeer.*$transcoderApiPort")
  if [ -n "${transcoderPid}" ]
  then
    transcoderRunning=true
  else
    transcoderRunning=false
  fi

}

function __lpdev_node_reset {

  pkill -9 livepeer

  unset broadcasterPid
  broadcasterRunning=false
  unset broadcasterGeth

  unset transcoderPid
  transcoderRunning=false
  unset transcoderGeth

}

function __lpdev_node_broadcaster {

  __lpdev_node_refresh_status

  if $broadcasterRunning
  then
    echo "Broadcaster running ($broadcasterPid)"
  fi

  __lpdev_geth_refresh_status
  __lpdev_protocol_refresh_status

  if ! $gethRunning || ! $protocolBuilt || [ -z ${controllerAddress} ]
  then
    echo "Geth must be running & protocol must be deployed to run a node"
    return 1
  fi

  if ! $broadcasterRunning
  then
    echo "Creating broadcaster account"
    broadcasterGeth=$(geth account new --password <(echo "") | cut -d' ' -f2 | tr -cd '[:alnum:]')
    echo "Created $broadcasterGeth"
  else
    broadcasterGeth=$(pgrep -fla "livepeer.*$broadcasterApiPort" | sed -nr 's/.*ethAccountAddr ([a-zA-Z0-9]+) .*/\1/p')
  fi

  if [ -z $broadcasterGeth ]
  then
    echo "Couldn't find the broadcast node's Account address"
    return 1
  fi

  echo "Transferring funds to $broadcasterGeth"
  transferEth="geth attach ipc:/home/ubuntu/.ethereum/geth.ipc --exec 'eth.sendTransaction({from: \"$gethMiningAccount\", to: \"$broadcasterGeth\", value: web3.toWei(1000000, \"ether\")})'"
  echo "Running $transferEth"
  eval $transferEth

  nodeDataDir=$nodeBaseDataDir/${broadcasterGeth:0:10}
  if [ ! -d $nodeDataDir ]
  then
    mkdir -p $nodeDataDir
  fi

  if ! $broadcasterRunning && [ -n $broadcasterGeth ]
  then
    echo "Running LivePeer broadcast node with the following command:
      $binDir -bootnode
              -controllerAddr $controllerAddress
              -datadir $nodeDataDir
              -ethDatadir $HOME/.ethereum
              -ethAccountAddr $broadcasterGeth
              -monitor=false
              -rtmp $broadcasterRtmpPort
              -http $broadcasterApiPort"

    $binDir -bootnode -controllerAddr $controllerAddress -datadir $nodeDataDir \
      -ethDatadir $HOME/.ethereum -ethAccountAddr $broadcasterGeth -monitor=false -rtmp $broadcasterRtmpPort \
      -http $broadcasterApiPort &>> $nodeDataDir/broadcaster.log &

    if [ $? -ne 0 ]
    then
      echo "Could not start LivePeer broadcast node"
      return 1
    else
      echo "LivePeer broadcast node started successfully"
      disown
      broadcasterRunning=true
    fi
  fi

  # Wait for the node's webserver to start
  echo -n "Attempting to connect to the LivePeer broadcast node webserver"
  attempts=15
  while ! nc -z localhost $broadcasterApiPort
  do
    if [ $attempts -eq 0 ]
    then
      echo "Giving up."
      return 1
    fi
    echo -n "."
    sleep 1
    attempts=$((attempts - 1))
  done

  echo ""

  echo "Requesting test tokens"
  curl -X "POST" http://localhost:$broadcasterApiPort/requestTokens

  echo "Depositing 500 tokens"
  curl -X "POST" http://localhost:$broadcasterApiPort/deposit \
    --data-urlencode "amount=500"

}

function __lpdev_node_transcoder {

  __lpdev_node_refresh_status

  if $transcoderRunning
  then
    echo "Transcoder running ($transcoderPid)"
  fi

  __lpdev_geth_refresh_status
  __lpdev_protocol_refresh_status

  if ! $gethRunning || ! $protocolBuilt || [ -z ${controllerAddress} ]
  then
    echo "Geth must be running & protocol must be deployed to run a node"
    return 1
  fi

  if ! $transcoderRunning
  then
    echo "Creating transcoder account"
    transcoderGeth=$(geth account new --password <(echo "") | cut -d' ' -f2 | tr -cd '[:alnum:]')
    echo "Created $transcoderGeth"
  else
    transcoderGeth=$(pgrep -fla "livepeer.*$transcoderApiPort" | sed -nr 's/.*ethAccountAddr ([a-zA-Z0-9]+) .*/\1/p')
  fi

  if [ -z $transcoderGeth ]
  then
    echo "Couldn't find the transcoder node's Account address"
    return 1
  fi

  echo "Transferring funds to $transcoderGeth"
  transferEth="geth attach ipc:/home/ubuntu/.ethereum/geth.ipc --exec 'eth.sendTransaction({from: \"$gethMiningAccount\", to: \"$transcoderGeth\", value: web3.toWei(1000000, \"ether\")})'"
  echo "Running $transferEth"
  eval $transferEth

  nodeDataDir=$nodeBaseDataDir/${transcoderGeth:0:10}
  if [ ! -d $nodeDataDir ]
  then
    mkdir -p $nodeDataDir
  fi

  bootNodePort=$(pgrep -fla "livepeer.*bootnode" | sed -nr "s/.*http ([0-9]+)( .*|$)/\1/p")
  if [ -n $bootNodePort ]
  then
    bootNodeId=$(curl http://localhost:$bootNodePort/nodeID 2> /dev/null)
    if [ -z $bootNodeId ]
    then
      echo "Could not find a boot node id (make sure you're running a node with the -bootnode flag)"
      return 1
    fi
  fi

  if ! $transcoderRunning && [ -n $transcoderGeth ]
  then
    echo "Running LivePeer transcode node with the following command:
      $binDir -controllerAddr $controllerAddress
              -datadir $nodeDataDir
              -ethDatadir $HOME/.ethereum
              -ethAccountAddr $transcoderGeth
              -monitor=false
              -rtmp $transcoderRtmpPort
              -http $transcoderApiPort
              -bootID $bootNodeId
              -bootAddr \"/ip4/localhost/tcp/15000\"
              -p 15001
              -transcoder"

    $binDir -p 15001 -controllerAddr $controllerAddress -datadir $nodeDataDir \
      -ethDatadir $HOME/.ethereum -ethAccountAddr $transcoderGeth -monitor=false -rtmp $transcoderRtmpPort \
      -http $transcoderApiPort -bootID $bootNodeId -bootAddr "/ip4/127.0.0.1/tcp/15000" \
      -transcoder &>> $nodeDataDir/transcoder.log &

    if [ $? -ne 0 ]
    then
      echo "Could not start LivePeer transcoder node"
      return 1
    else
      echo "LivePeer transcoder node started successfully"
      disown
      transcoderRunning=true
    fi
  fi

  # Wait for the node's webserver to start
  echo -n "Attempting to connect to the LivePeer transcoder node webserver"
  attempts=15
  while ! nc -z localhost $transcoderApiPort
  do
    if [ $attempts -eq 0 ]
    then
      echo "Giving up."
      return 1
    fi
    echo -n "."
    sleep 1
    attempts=$((attempts - 1))
  done

  echo ""

  echo "Requesting test tokens"
  curl -X "POST" http://localhost:$transcoderApiPort/requestTokens

  echo "Activating transcoder"
  curl -X "POST" http://localhost:$transcoderApiPort/deposit \
    --data-urlencode "blockRewardCut=10&feeShare=5&pricePerSegment=1&amount=500"

}

function __lpdev_wizard {

  echo "
+----------------------------------------------------+
| Welcome to the Livepeer local dev environment tool |
|                                                    |
+----------------------------------------------------+
"
  __lpdev_status

  echo "What would you like to do?"

  wizardOptions=(
  "Display status"
  "Set up & start Geth local network"
  "Deploy/overwrite protocol contracts"
  #"Set up IPFS"
  "Start & set up broadcaster node"
  "Start & set up transcoder node"
  #"Deposit tokens to node"
  "Destroy current environment"
  "Exit"
  )

  select opt in "${wizardOptions[@]}"
  do
    case $opt in
      "Display status")
        __lpdev_status
        ;;
      "Set up & start Geth local network")
        __lpdev_geth_init
        __lpdev_geth_run
        ;;
      "Deploy/overwrite protocol contracts")
        __lpdev_protocol_init
        __lpdev_protocol_deploy
        ;;
      "Set up IPFS")
        echo "Coming soon";;
      "Start & set up broadcaster node")
        __lpdev_node_broadcaster
        ;;
      "Start & set up transcoder node")
        __lpdev_node_transcoder
        ;;
      "Deposit tokens to node")
        echo "Coming soon";;
      "Destroy current environment")
        __lpdev_reset
        ;;
      "Exit")
        return 0;;
    esac
  done
}

alias lpdev=__lpdev_wizard
