package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	erpc "github.com/Varunram/essentials/rpc"
)

func checkReq(w http.ResponseWriter, r *http.Request) ([]string, error) {
	var arr []string
	err := erpc.CheckPost(w, r) // check origin of request as well if needed
	if err != nil {
		log.Println(err)
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return arr, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return arr, err
	}
	var rf RequestFormat
	err = json.Unmarshal(data, &rf)
	if err != nil {
		log.Println(err)
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		return arr, err
	}

	arr = rf.Addresses
	return arr, nil
}

// MultigetUtxos gets the utxos associated with multiple addresses
func MultigetUtxos() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigetutxos", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

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

func multiAddr(w http.ResponseWriter, r *http.Request,
	arr []string) ([]MultigetAddrReturn, error) {

	x := make([]MultigetAddrReturn, len(arr))
	currentBh, err := CurrentBlockHeight()
	if err != nil {
		log.Println(err)
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		return x, err
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
	return x, nil
}

// MultigetAddr gets all data associated with a particular address
func MultigetAddr() {
	// make a curl request out to localhost and get the ping response
	http.HandleFunc("/multiaddr", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		x, err := multiAddr(w, r, arr)
		if err != nil {
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// GetBalAndTx combines the balance and Multigetaddr endpoints
func GetBalAndTx() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/baltxs", func(w http.ResponseWriter, r *http.Request) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var ret BalTx
		ret.Balance = multiBalance(arr, w, r)
		ret.Transactions, err = multiAddr(w, r, arr)
		if err != nil {
			return
		}

		erpc.MarshalSend(w, ret)
	})
}

func multiBalance(arr []string, w http.ResponseWriter, r *http.Request) MultigetBalanceReturn {
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
	return x
}

// MultigetBalance gets the net balance associated with multiple addresses
func MultigetBalance() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigetbalance", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var x MultigetBalanceReturn
		x = multiBalance(arr, w, r)

		erpc.MarshalSend(w, x)
	})
}

// MultigetTxs gets the transactions associated with mutliple addresses
func MultigetTxs() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigettxs", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

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
