package main

import (
	"log"
	"net/http"

	//"strings"

	erpc "github.com/Varunram/essentials/rpc"
)

func startHandlers() {
	MultigetAddr()
	GetBalAndTx()
	MultigetUtxos()
	MultigetBalance()
	MultigetTxs()

	erpc.SetupPingHandler()
	GetFees()
	PostTx()
	RelayTxid()
	RelayGetRequest()
}

func main() {
	startHandlers()
	// if you're running esplora, use socat tcp-listen:3003,reuseaddr,fork tcp:localhost:3002 to tunnel port since
	// it does not seem possible to open the port directly
	// // setup https here
	err := http.ListenAndServeTLS("localhost:445", "ssl/server.crt", "ssl/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
