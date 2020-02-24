package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	erpc "github.com/Varunram/essentials/rpc"
)

// MultigetUtxos gets the utxos associated with multiple addresses
func MultigetUtxos() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigetutxos", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		arr := rf.Addresses
		var result [][]Utxo
		for _, elem := range arr {
			// send the request out
			go func(elem string) {
				tempTxs, err := GetUtxosAddress(w, r, elem)
				if err != nil {
					log.Println(err)
					erpc.ResponseHandler(w, http.StatusInternalServerError)
					return
				}
				result = append(result, tempTxs)
			}(elem)
		}
		time.Sleep(50 * time.Millisecond)
		erpc.MarshalSend(w, result)
	})
}

// MultigetAddr gets all data associated with a particular address
func MultigetAddr() {
	// make a curl request out to localhost and get the ping response
	http.HandleFunc("/multiaddr", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		arr := rf.Addresses
		x := make([]MultigetAddrReturn, len(arr))
		currentBh, err := CurrentBlockHeight()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		for i, elem := range arr {
			x[i].Address = elem // store the address of the passed elements
			// send the request out
			go func(i int, elem string) {
				allTxs, err := GetTxsAddress(w, r, elem)
				if err == nil {
					x[i].TotalTransactions = float64(len(allTxs))
					x[i].Transactions = allTxs
					for j := range x[i].Transactions {
						if x[i].Transactions[j].Status.Confirmed {
							x[i].Transactions[j].NumberofConfirmations = currentBh - x[i].Transactions[j].Status.BlockHeight
						} else {
							x[i].Transactions[j].NumberofConfirmations = 0
						}
					}
					x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = 0, 0
					go func(i int, elem string) {
						x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = GetBalanceCount(w, r, elem)
					}(i, elem)
				}
			}(i, elem)
		}
		time.Sleep(50 * time.Millisecond)
		erpc.MarshalSend(w, x)
	})
}

// GetBalAndTx combines the balance and Multigetaddr endpoints
func GetBalAndTx() {

	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/baltxs", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("unable to read body: ", err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		arr := rf.Addresses

		var ret BalTx

		for _, elem := range arr {
			// send the request out
			tBalance, tUnconfirmedBalance := 0.0, 0.0
			go func(elem string) {
				tBalance, tUnconfirmedBalance = GetBalanceAddress(w, r, elem)
				ret.Balance.Balance += tBalance
				ret.Balance.UnconfirmedBalance += tUnconfirmedBalance
			}(elem)
		}

		// the following call is blocking, so we don't need to add a sleep routine for the
		// last call above to work
		x := make([]MultigetAddrReturn, len(arr))
		currentBh, err := CurrentBlockHeight()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		for i, elem := range arr {
			x[i].Address = elem // store the address of the passed elements
			// send the request out
			go func(i int, elem string) {
				allTxs, err := GetTxsAddress(w, r, elem)
				if err == nil {
					x[i].TotalTransactions = float64(len(allTxs))
					x[i].Transactions = allTxs
					for j := range x[i].Transactions {
						if x[i].Transactions[j].Status.Confirmed {
							x[i].Transactions[j].NumberofConfirmations = currentBh - x[i].Transactions[j].Status.BlockHeight
						} else {
							x[i].Transactions[j].NumberofConfirmations = 0
						}
					}
					x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = 0, 0
					go func(i int, elem string) {
						x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = GetBalanceCount(w, r, elem)
					}(i, elem)
				}
			}(i, elem)
		}
		time.Sleep(50 * time.Millisecond)
		ret.Transactions = x
		erpc.MarshalSend(w, ret)
	})
}

// MultigetBalance gets the net balance associated with multiple addresses
func MultigetBalance() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigetbalance", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		arr := rf.Addresses
		var x MultigetBalanceReturn

		for _, elem := range arr {
			// send the request out
			tBalance, tUnconfirmedBalance := 0.0, 0.0
			go func(elem string) {
				tBalance, tUnconfirmedBalance = GetBalanceAddress(w, r, elem)
				x.Balance += tBalance
				x.UnconfirmedBalance += tUnconfirmedBalance
			}(elem)
		}

		time.Sleep(50 * time.Millisecond)
		erpc.MarshalSend(w, x)
	})
}

// MultigetTxs gets the transactions associated with mutliple addresses
func MultigetTxs() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigettxs", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		arr := rf.Addresses
		var x TxReturn
		for _, elem := range arr {
			// send the request out
			go func(elem string) {
				tempTxs, err := GetTxsAddress(w, r, elem)
				if err != nil {
					log.Println(err)
					erpc.ResponseHandler(w, http.StatusInternalServerError)
					return
				}
				x.Txs = append(x.Txs, tempTxs)
			}(elem)
		}

		time.Sleep(50 * time.Millisecond)
		erpc.MarshalSend(w, x)
	})
}
