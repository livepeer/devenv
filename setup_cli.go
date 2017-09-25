package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/glog"
	cli "gopkg.in/urfave/cli.v1"
)

var dir string
var broadcasterPort = 8935
var transcoderPort = 8936

func main() {
	app := cli.NewApp()
	app.Name = "testenv-cli"
	app.Usage = "Set up local Livepeer testing environment"

	app.Action = func(c *cli.Context) error {
		// Set up the logger to print everything and the random generator
		log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(c.Int("loglevel")), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
		rand.Seed(time.Now().UnixNano())

		w := &wizard{
			in: bufio.NewReader(os.Stdin),
		}
		w.run()
		return nil
	}

	var err error
	dir, err = os.Getwd()
	if err != nil {
		glog.Errorf("Error getting wd: %v", err)
		return
	}

	app.Run(os.Args)
}

type wizard struct {
	in           *bufio.Reader // Wrapper around stdin to allow reading user input
	ProtocolAddr string
	TokenAddr    string
	FaucetAddr   string
}

func (w *wizard) run() {
	glog.Infof("Make sure you are in the testenv directory.")

	fmt.Println("+-------------------------------------------------------------+")
	fmt.Println("| Welcome to testenv-cli, your Livepeer test environment tool |")
	fmt.Println("|                                                             |")
	fmt.Println("+-------------------------------------------------------------+")
	fmt.Println()

	for {
		w.stats()
		fmt.Println()
		fmt.Println("What would you like to do?")
		fmt.Println(" 1. Set up seed data")
		fmt.Println(" 2. Set up & start Geth")
		fmt.Println(" 3. Deploy new protocol contracts")
		fmt.Println(" 4. Set up IPFS")
		fmt.Println(" 5. Start & Set up broadcaster node")
		fmt.Println(" 6. Start & Set up transcoder node")

		choice := w.read()
		switch {
		case choice == "1":
			w.setupSeedData()
		case choice == "2":
			w.setupAndStartGeth()
		case choice == "3":
			w.deployProtocol()
		case choice == "4":
			glog.Infof("TODO...")
		case choice == "5":
			w.setupAndStartBroadcaster()
		case choice == "6":
			w.setupAndStartTranscoder()
		default:
			log.Error("That's not something I can do")
		}
	}
}
