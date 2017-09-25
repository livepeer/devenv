package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
		// fmt.Println(err.Error())
		return []byte{}
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
	tc := string(File(filepath.Join(dir, "protocol/build/contracts/LivepeerToken.json")))
	if tc != "" {
		w.TokenAddr = strings.Split(tc, "address\": \"")[1][:42]
	}
	fmt.Fprintf(wtr, "LivepeerToken: \t%v\n", w.TokenAddr)

	pc := string(File(filepath.Join(dir, "protocol/build/contracts/LivepeerProtocol.json")))
	if pc != "" {
		w.ProtocolAddr = strings.Split(pc, "address\": \"")[1][:42]
	}
	fmt.Fprintf(wtr, "LivepeerProtocol: \t%v\n", w.ProtocolAddr)

	fc := string(File(filepath.Join(dir, "protocol/build/contracts/LivepeerTokenFaucet.json")))
	if fc != "" {
		w.FaucetAddr = strings.Split(fc, "address\": \"")[1][:42]
	}
	fmt.Fprintf(wtr, "LivepeerTokenFaucet: \t%v\n", w.FaucetAddr)

	//TODO: Print out if IPFS is set up and running

	//Print out if broadcaster is running
	bStatus := []string{}
	if _, err := http.Get(fmt.Sprintf("http://localhost:%v", broadcasterPort)); err == nil {
		bStatus = append(bStatus, "Running")
	} else {
		bStatus = append(bStatus, "Not Running")
	}
	if w.getDeposit(broadcasterPort) == 0 {
		bStatus = append(bStatus, "No Deposit")
	} else {
		bStatus = append(bStatus, "Has Deposit")
	}
	fmt.Fprintf(wtr, "Broadcaster running: \t%v\n", strings.Join(bStatus, ", "))

	//Print out if transcoder is running
	tStatus := []string{}
	if _, err := http.Get(fmt.Sprintf("http://localhost:%v", transcoderPort)); err == nil {
		tStatus = append(tStatus, "Running")
	} else {
		tStatus = append(tStatus, "Not Running")
	}
	tStatus = append(tStatus, w.getTranscoderStatus(transcoderPort))
	fmt.Fprintf(wtr, "Transcoder running: \t%v\n", strings.Join(tStatus, ", "))

	wtr.Flush()
}

func (w *wizard) setupSeedData() {
	//Set up lpdata dirs
	// if _, err := os.Stat(filepath.Join(dir, "lpdata1")); os.IsNotExist(err) {
	os.Mkdir("lpdata1", 0777)
	Copy(filepath.Join(dir, "lpkeys1.json"), filepath.Join(dir, "lpdata1", "keys.json"))
	os.Mkdir("lpdata2", 0777)
	Copy(filepath.Join(dir, "lpkeys2.json"), filepath.Join(dir, "lpdata2", "keys.json"))
	// }
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
	glog.Infof("Run `npm install && truffle migrate --network lpTestNet` in %v", filepath.Join(dir, "protocol"))

	//TODO: Run truffle migrate --network lpTestNet
}

func (w *wizard) setupAndStartBroadcaster() {
	if _, err := os.Stat(filepath.Join(dir, "livepeer")); os.IsNotExist(err) {
		//Download Livepeer
		download("https://github.com/livepeer/go-livepeer/releases/download/0.1.3/livepeer_darwin", "livepeer")
		//TODO: Change to tar deployment and untar the executables.
	}
	cmd := fmt.Sprintf("./livepeer -bootnode -protocolAddr %v -tokenAddr %v -faucetAddr %v -datadir %v -ethDatadir %v -ethAccountAddr 0x94107cb2261e722f9f4908115546eeee17decada -monitor=false -rtmp %v -http %v &> lpBroadcaster.log",
		w.ProtocolAddr, w.TokenAddr, w.FaucetAddr, filepath.Join(dir, "lpdata1"), filepath.Join(dir, "lpGeth"), broadcasterPort-7000, broadcasterPort)
	glog.Infof("Command: %v", cmd)

	if w.getTokenBalance(broadcasterPort) == 0 {
		glog.Infof("Requesting for test tokens")
		httpPost(fmt.Sprintf("http://localhost:%v/requestTokens", broadcasterPort))
	}

	if w.getDeposit(broadcasterPort) == 0 {
		glog.Infof("Depositing 500 tokens")
		val := url.Values{
			"amount": {"500"},
		}
		httpPostWithParams(fmt.Sprintf("http://localhost:%v/deposit", broadcasterPort), val)
	}

}

func (w *wizard) setupAndStartTranscoder() {
	if _, err := os.Stat(filepath.Join(dir, "livepeer")); os.IsNotExist(err) {
		//Download Livepeer
		download("https://github.com/livepeer/go-livepeer/releases/download/0.1.3/livepeer_darwin", "livepeer")
		//TODO: Change to tar deployment and untar the executables.
	}
	cmd := fmt.Sprintf("./livepeer -bootnode -protocolAddr %v -tokenAddr %v -faucetAddr %v -datadir %v -ethDatadir %v -ethAccountAddr 0x94107cb2261e722f9f4908115546eeee17decada -monitor=false -rtmp %v -http %v -bootID 1220fd39d26e8a0ddf25693e574b821df32c45cd18fee1ee8bf329da96eb67bd2b5f -bootAddr /ip4/127.0.0.1/tcp/15000 -p 15001 -transcoder &> lpBroadcaster.log",
		w.ProtocolAddr, w.TokenAddr, w.FaucetAddr, filepath.Join(dir, "lpdata2"), filepath.Join(dir, "lpGeth"), transcoderPort-7000, transcoderPort)
	glog.Infof("Command: %v", cmd)

	if w.getTranscoderStatus(transcoderPort) != "Active" {
		glog.Infof("Activating transcoder")

		val := url.Values{
			"blockRewardCut":  {fmt.Sprintf("%v", 10)},
			"feeShare":        {fmt.Sprintf("%v", 5)},
			"pricePerSegment": {fmt.Sprintf("%v", 1)},
			"amount":          {fmt.Sprintf("%v", 500)},
		}

		httpPostWithParams(fmt.Sprintf("http://localhost:%v/activateTranscoder", transcoderPort), val)
	}
}

func download(rawURL, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(rawURL) // add a filter to check redirect

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s with %v bytes downloaded", fileName, size)
}

func (w *wizard) getDeposit(port int) int {
	e := httpGet(fmt.Sprintf("http://localhost:%v/broadcasterDeposit", port))
	i, _ := strconv.Atoi(e)
	return i
}

func (w *wizard) getTokenBalance(port int) int {
	b := httpGet(fmt.Sprintf("http://localhost:%v/tokenBalance", port))
	i, _ := strconv.Atoi(b)
	return i
}

func (w *wizard) getTranscoderStatus(port int) string {
	return httpGet(fmt.Sprintf("http://localhost:%v/transcoderStatus", port))
}
