package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bithyve/bithyve-wrapper/electrs"

	//"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"
)

var opts struct {
	Mainnet bool `short:"m" description:"Connect to mainnet"`
}

func startHandlers() {
	router := httprouter.New()
	MultiData(router)
	MultiBalTxs(router)
	MultiUtxos(router)
	MultiBalances(router)
	MultiTxs(router)

	GetFees(router)
	PostTx(router)
	RelayTxid()
	RelayGetRequest()

	ping(router)

	err := http.ListenAndServeTLS("localhost:445", "ssl/server.crt", "ssl/server.key", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	// if you're running esplora, use socat tcp-listen:3003,reuseaddr,fork tcp:localhost:3002 to tunnel port since
	// it does not seem possible to open the port directly
	// // setup https here

	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Mainnet {
		log.Println("connecting to electrs mainnet")
		electrs.SetMainnet()
	}

	startHandlers()
}
