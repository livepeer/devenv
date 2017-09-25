package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/golang/glog"
)

type Addresses struct {
	Protocol string `json:"protocol"`
	Token    string `json:"token"`
	Faucet   string `json:"faucet"`
}

func setup() {
	dir, err := os.Getwd()
	if err != nil {
		glog.Errorf("Error: %v", err)
		return
	}

	// //Set up geth accounts
	// if _, err := os.Stat(filepath.Join(dir, "geth1")); os.IsNotExist(err) {
	// 	glog.Infof("Setting up eth account 1 in %v", filepath.Join(dir, "geth1"))
	// 	cmd1 := exec.Command("geth", "-datadir", "./geth1", "account", "new")
	// 	stdin, err := cmd1.StdinPipe()
	// 	if err != nil {
	// 		glog.Errorf("Error: %v", err)
	// 	}
	// 	stdin.Write([]byte("\n\n\n"))
	// 	cmd1.Start()
	// 	go func() {
	// 		if err := cmd1.Wait(); err != nil {
	// 			glog.Infof("Couldn't start geth: %v", err)
	// 			os.Exit(1)
	// 		}
	// 	}()
	// }

	// if _, err := os.Stat(filepath.Join(dir, "geth2")); os.IsNotExist(err) {
	// 	glog.Infof("Setting up eth account 2 in %v", filepath.Join(dir, "geth2"))
	// 	cmd2 := exec.Command("geth", "-datadir", "./geth2", "account", "new")
	// 	stdin, err := cmd2.StdinPipe()
	// 	if err != nil {
	// 		glog.Errorf("Error: %v", err)
	// 	}
	// 	stdin.Write([]byte("\n\n\n"))
	// 	cmd2.Start()
	// 	go func() {
	// 		if err := cmd2.Wait(); err != nil {
	// 			glog.Infof("Couldn't start geth: %v", err)
	// 			os.Exit(1)
	// 		}
	// 	}()
	// }

	//Set up lpdata dirs
	if _, err := os.Stat(filepath.Join(dir, "lpdata1")); os.IsNotExist(err) {
		os.Mkdir("lpdata1", 777)
		Copy(filepath.Join(dir, "lpkeys1.json"), filepath.Join(dir, "lpdata1", "keys.json"))
		os.Mkdir("lpdata2", 777)
		Copy(filepath.Join(dir, "lpkeys2.json"), filepath.Join(dir, "lpdata2", "keys.json"))
	}

	//Set up geth genesis, init geth data
	gethdir := filepath.Join(dir, "lpGeth")
	if _, err := os.Stat(gethdir); os.IsNotExist(err) {
		glog.Infof("Setting up geth data at %v", gethdir)
		cmd := exec.Command("geth", "-datadir", gethdir, "init", "localPOA.json")
		if err := cmd.Start(); err != nil {
			glog.Fatalf("Error: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			glog.Errorf("Error: %v", err)
		}
		//Copy keys into lpGeth
		Copy(filepath.Join(dir, "keystore", "key1"), filepath.Join(dir, "lpGeth", "keystore", "key1"))
		Copy(filepath.Join(dir, "keystore", "key2"), filepath.Join(dir, "lpGeth", "keystore", "key2"))
	}
	glog.Infof("To start Geth, run: `geth -datadir %v -networkid 54321 -rpc -unlock 0x94107cb2261e722f9f4908115546eeee17decada -mine console`", gethdir)

	//Set up protocol contract
	if _, err := os.Stat(filepath.Join(dir, "protocol")); os.IsNotExist(err) {
		glog.Infof("Setting up protocol")
		//Clone the protocol repo
		cmd := exec.Command("git", "clone", "git@github.com:livepeer/protocol.git")
		cmd.Start()
		if err := cmd.Wait(); err != nil {
			glog.Infof("Couldn't clone: %v", err)
			os.Exit(1)
		}
		glog.Infof("Done cloning protocol")

		//Copy truffle.js
		Copy(filepath.Join(dir, "truffle.js"), filepath.Join(dir, "protocol", "truffle.js"))
		Copy(filepath.Join(dir, "migrations.config.js"), filepath.Join(dir, "protocol", "migrations", "truffle.js"))
		glog.Infof("Run `npm install` in %v", filepath.Join(dir, "protocol"))
	}
	glog.Infof("To deploy contracts, run `truffle migrate --network lpTestNet`")

	datadir1 := filepath.Join(dir, "lpdata1")
	ethdatadir := gethdir
	contracts := Contracts(filepath.Join(dir, "contracts.json"))
	glog.Infof("Contracts: %v", contracts)

	//Start 2 livepeer nodes
	lpcmd1 := exec.Command("livepeer", "-testnet", "-bootnode",
		"-protocolAddr", contracts.Protocol, "-tokenAddr", contracts.Token, "-faucetAddr", contracts.Faucet,
		"-datadir", datadir1, "-ethDatadir", ethdatadir, "-ethAccountAddr", "0x94107cb2261e722f9f4908115546eeee17decada")
	if err := lpcmd1.Start(); err != nil {
		glog.Errorf("Error starting livepeer node: %v", err)
	}

	datadir2 := filepath.Join(dir, "lpdata2")
	lpcmd2 := exec.Command("livepeer", "-testnet",
		"-protocolAddr", contracts.Protocol, "-tokenAddr", contracts.Token, "-faucetAddr", contracts.Faucet,
		"-datadir", datadir2, "-ethDatadir", ethdatadir, "-ethAccountAddr", "0x0ddb225031ccb58ff42866f82d907f7766899014",
		"-rtmp", "1936", "-http", "8936", "-bootID", "1220fd39d26e8a0ddf25693e574b821df32c45cd18fee1ee8bf329da96eb67bd2b5f", "-bootAddr", "/ip4/127.0.0.1/tcp/15000", "-p", "15001", "-transcoder")
	if err := lpcmd1.Start(); err != nil {
		glog.Errorf("Error starting livepeer node: %v", err)
	}

	if err := lpcmd2.Wait(); err != nil {
		glog.Fatalf("Error: %v", err)
	}

	//Set up
}

func Copy(src string, dst string) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		glog.Fatalf("Error reading: %v", err)
	}

	f, _ := os.Create(dst)
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		glog.Fatalf("Error writing: %v", err)
	}
}

func File(filename string) []byte {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return raw
}

func Contracts(filename string) Addresses {
	var c Addresses
	json.Unmarshal(File(filename), &c)
	return c
}

func (w *wizard) stats() {
	wtr := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)

	//Print out if seed data is set up
	seedSetup := false
	if _, err := os.Stat(filepath.Join(dir, "lpdata1/keys.json")); !os.IsNotExist(err) {
		if _, err := os.Stat(filepath.Join(dir, "lpdata2/keys.json")); !os.IsNotExist(err) {
			seedSetup = true
		}
	}
	fmt.Fprintf(wtr, "Seed data setup: \t%v\n", seedSetup)

	//Print out if geth is set up and running
	gethSetup := false
	gethdir := filepath.Join(dir, "lpGeth")
	if _, err := os.Stat(filepath.Join(gethdir, "geth.ipc")); !os.IsNotExist(err) {
		gethSetup = true
	}
	fmt.Fprintf(wtr, "Geth setup and running: \t%v\n", gethSetup)

	//Print out protocol addresses
	token := ""
	tc := string(File(filepath.Join(dir, "protocol/build/contracts/LivepeerToken.json")))
	if tc != "" {
		token = strings.Split(tc, "address\": \"")[1][:42]
	}
	fmt.Fprintf(wtr, "LivepeerToken: \t%v\n", token)

	protocol := ""
	pc := string(File(filepath.Join(dir, "protocol/build/contracts/LivepeerProtocol.json")))
	if pc != "" {
		protocol = strings.Split(pc, "address\": \"")[1][:42]
	}
	fmt.Fprintf(wtr, "LivepeerProtocol: \t%v\n", protocol)

	faucet := ""
	fc := string(File(filepath.Join(dir, "protocol/build/contracts/LivepeerTokenFaucet.json")))
	if fc != "" {
		faucet = strings.Split(fc, "address\": \"")[1][:42]
	}
	fmt.Fprintf(wtr, "LivepeerTokenFaucet: \t%v\n", faucet)

	//TODO: Print out if IPFS is set up and running

	//Print out if broadcaster is running
	broadcaster := true
	if _, err := http.Get(fmt.Sprintf("http://localhost:%v", broadcasterPort)); err != nil {
		broadcaster = false
	}
	fmt.Fprintf(wtr, "Broadcaster running: \t%v\n", broadcaster)

	//Print out if transcoder is running
	transcoder := true
	if _, err := http.Get(fmt.Sprintf("http://localhost:%v", transcoderPort)); err != nil {
		transcoder = false
	}
	fmt.Fprintf(wtr, "Transcoder running: \t%v\n", transcoder)

	wtr.Flush()
}

func (w *wizard) setupSeedData() {
	//Set up lpdata dirs
	if _, err := os.Stat(filepath.Join(dir, "lpdata1")); os.IsNotExist(err) {
		os.Mkdir("lpdata1", 777)
		Copy(filepath.Join(dir, "lpkeys1.json"), filepath.Join(dir, "lpdata1", "keys.json"))
		os.Mkdir("lpdata2", 777)
		Copy(filepath.Join(dir, "lpkeys2.json"), filepath.Join(dir, "lpdata2", "keys.json"))
	}
}

func (w *wizard) setupAndStartGeth() {
	//Set up geth genesis, init geth data
	gethdir := filepath.Join(dir, "lpGeth")
	if _, err := os.Stat(gethdir); os.IsNotExist(err) {
		glog.Infof("Setting up geth data at %v", gethdir)
		cmd := exec.Command("geth", "-datadir", gethdir, "init", "localPOA.json")
		if err := cmd.Start(); err != nil {
			glog.Fatalf("Error: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			glog.Errorf("Error: %v", err)
		}
		//Copy keys into lpGeth
		Copy(filepath.Join(dir, "keystore", "key1"), filepath.Join(dir, "lpGeth", "keystore", "key1"))
		Copy(filepath.Join(dir, "keystore", "key2"), filepath.Join(dir, "lpGeth", "keystore", "key2"))
	}
	glog.Infof("To start Geth, run: `geth -datadir %v -networkid 54321 -rpc -unlock 0x94107cb2261e722f9f4908115546eeee17decada -mine console`", gethdir)

}

func (w *wizard) deployProtocol() {
	//Set up protocol contract
	if _, err := os.Stat(filepath.Join(dir, "protocol/build/contracts/LivepeerToken")); os.IsExist(err) {
		glog.Infof("Already deployed protocol.  Want to refresh? (Y/n)")
		yn := w.read()
		if strings.ToLower(yn) == "n" {
			return
		}
	}

	glog.Infof("Setting up protocol")
	//Clone the protocol repo
	if _, err := os.Stat(filepath.Join(dir, "protocol")); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "git@github.com:livepeer/protocol.git")
		cmd.Start()
		if err := cmd.Wait(); err != nil {
			glog.Infof("Couldn't clone: %v", err)
			os.Exit(1)
		}
		glog.Infof("Done cloning protocol")

		//Copy truffle.js
		Copy(filepath.Join(dir, "truffle.js"), filepath.Join(dir, "protocol", "truffle.js"))
		Copy(filepath.Join(dir, "migrations.config.js"), filepath.Join(dir, "protocol", "migrations", "truffle.js"))
	} else {
		// exec.Command("cd", "protocol", "&&", "git", "pull")
		// cmd.Start()
		// if err := cmd.Wait(); err != nil {
		// 	glog.Infof("Couldn't clone: %v", err)
		// 	os.Exit(1)
		// }
	}
	glog.Infof("Run `npm install` in %v", filepath.Join(dir, "protocol"))

	//TODO: Run truffle migrate --network lpTestNet
	glog.Infof("To deploy contracts, run `truffle migrate --network lpTestNet`")

	//TODO: Parse and print out contract addresses

}

func (w *wizard) setupBroadcaster() {

}
