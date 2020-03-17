package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bithyve/bithyve-wrapper/format"

	erpc "github.com/Varunram/essentials/rpc"
	electrs "github.com/bithyve/bithyve-wrapper/electrs"
)

// RelayTxid gets the information associated with a particular tx on the blockchain
func RelayTxid() {
	http.HandleFunc("/txid", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		if r.URL.Query()["txid"] == nil {
			erpc.ResponseHandler(w, http.StatusBadRequest, "required param txid not found")
			return
		}

		txid := r.URL.Query()["txid"][0]
		body := electrs.ElectrsURL + "/tx/" + txid
		var x format.Tx
		erpc.GetAndSendJson(w, body, x)
	})
}

// RelayGetRequest relays all remaining get requests to the esplora instance
func RelayGetRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		body := electrs.ElectrsURL + "" + r.URL.String()
		data, err := erpc.GetRequest(body)
		if err != nil {
			erpc.ResponseHandler(w, http.StatusInternalServerError, "error while relaying get request to electrs")
			log.Println("could not submit transacation to testnet, quitting")
			return
		}
		var x interface{}
		_ = json.Unmarshal(data, &x)
		erpc.MarshalSend(w, x)
	})
}
