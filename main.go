package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/bithyve/bithyve-wrapper/electrs"

	//"strings"

	erpc "github.com/Varunram/essentials/rpc"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	DevEnv     bool   `short:"d" description:"Start dev env"`
	Mainnet    bool   `short:"m" description:"Connect to mainnet"`
	ElectrsURL string `short:"u" description:"Connect to your own electrs instance"`
	BackupURL  string `short:"b" description:"Connect to your own backup electrs instance"`
	Test       bool   `short:"t" description:"Use for testing"`
	Logs       bool   `short:"l" description:"Testing logs"`
}

func startHandlers() {
	MultiData()
	MultiBalTxs()
	MultiUtxoTxs()
	MultiUtxos()
	MultiBalances()
	MultiTxs()
	NewMultiUtxoTxs()

	erpc.SetupPingHandler()
	GetFees(opts.Mainnet)
	GetFeesE(opts.Mainnet)
	PostTx()
	RelayTxid()
	RelayGetRequest()
}

func main() {

	var err error
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	startHandlers()

	if opts.Mainnet {
		log.Println("connecting to electrs mainnet")
		electrs.SetMainnet()
	}

	if opts.Logs {
		log.Println("starting in logs mode")
		electrs.ToggleLogs()
	}

	if opts.DevEnv {
		log.Println("starting in devenv mode")
		electrs.SetDevEnv()
	}

	if opts.ElectrsURL != "" {
		if opts.BackupURL != "" {
			electrs.SetURL(opts.ElectrsURL, opts.BackupURL)
		} else {
			if opts.Mainnet {
				electrs.SetURL(opts.ElectrsURL, "https://api.bithyve.com")
			} else {
				electrs.SetURL(opts.ElectrsURL, "https://test-wrapper.bithyve.com")
			}
		}
	}

	if opts.Test {
		err = http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		err = http.ListenAndServeTLS("localhost:445", "ssl/server.crt", "ssl/server.key", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
