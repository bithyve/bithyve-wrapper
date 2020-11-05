package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bithyve/bithyve-wrapper/electrs"

	//"strings"

	erpc "github.com/Varunram/essentials/rpc"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Mainnet bool `short:"m" description:"Connect to mainnet"`
	Test    bool `short:"t" description:"Use for testing"`
	Logs    bool `short:"l" description:"Testing logs"`
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

	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	startHandlers()

	if opts.Mainnet {
		log.Println("connecting to electrs mainnet")
		electrs.SetMainnet()
	}

	if opts.Logs {
		electrs.ToggleLogs()
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
