package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bithyve/bithyve-wrapper/electrs"
	"github.com/bithyve/bithyve-wrapper/format"
	"github.com/julienschmidt/httprouter"

	erpc "github.com/Varunram/essentials/rpc"
)

func wait() {
	time.Sleep(100 * time.Millisecond)
}

func blockWait(length int) {
	if length < 5 {
		time.Sleep(40 * time.Millisecond)
	} else if length >= 5 && length < 10 {
		time.Sleep(80 * time.Millisecond)
	} else if length >= 10 && length < 100 {
		time.Sleep(120 * time.Millisecond)
	} else if length >= 100 && length < 150 {
		time.Sleep(150 * time.Millisecond)
	} else if length >= 150 && length < 200 {
		time.Sleep(200 * time.Millisecond)
	} else {
		time.Sleep(500 * time.Millisecond)
	}
}

func checkReq(w http.ResponseWriter, r *http.Request) ([]string, error) {
	var arr []string
	err := erpc.CheckPost(w, r)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusNotFound)
		log.Println(err)
		return arr, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		log.Println(err)
		return arr, err
	}
	var rf format.RequestFormat
	err = json.Unmarshal(data, &rf)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		log.Println(err)
		return arr, err
	}

	arr = rf.Addresses

	// filter through list to remove duplicates
	nodupsMap := make(map[string]bool)
	var nodups []string

	for _, elem := range arr {
		if _, value := nodupsMap[elem]; !value {
			nodupsMap[elem] = true
			nodups = append(nodups, elem)
		}
	}

	return nodups, nil
}

func multiAddr(w http.ResponseWriter, r *http.Request,
	arr []string) ([]format.MultigetAddrReturn, error) {

	x := make([]format.MultigetAddrReturn, len(arr))
	currentBh, err := electrs.CurrentBlockHeight()
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		log.Println(err)
		return x, err
	}

	var maxTxs = 0
	if opts.Mainnet {
		for i, elem := range arr {
			x[i].Address = elem // store the address of the passed elements
			allTxs, err := electrs.GetTxsAddress(elem)
			if err == nil {
				if len(allTxs) > maxTxs {
					maxTxs = len(allTxs)
				}
				go func(i int, elem string, allTxs []format.Tx) {
					x[i].TotalTransactions = float64(len(allTxs))
					x[i].Transactions = allTxs
					x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = 0, 0
					for j := range x[i].Transactions {
						if x[i].Transactions[j].Status.Confirmed {
							x[i].Transactions[j].NumberofConfirmations =
								currentBh - x[i].Transactions[j].Status.BlockHeight
						} else {
							x[i].Transactions[j].NumberofConfirmations = 0
						}
					}
					go func(i int, elem string) {
						x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions =
							electrs.GetBalanceCount(elem)
					}(i, elem)
				}(i, elem, allTxs)
			} else {
				log.Println("error in gettxsaddress call: ", err)
			}
		}
		blockWait(maxTxs)
	} else {
		for i, elem := range arr {
			x[i].Address = elem // store the address of the passed elements
			allTxs, err := electrs.GetTxsAddress(elem)
			if err == nil {

				x[i].TotalTransactions = float64(len(allTxs))
				x[i].Transactions = allTxs
				x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = 0, 0
				for j := range x[i].Transactions {
					if x[i].Transactions[j].Status.Confirmed {
						x[i].Transactions[j].NumberofConfirmations =
							currentBh - x[i].Transactions[j].Status.BlockHeight
					} else {
						x[i].Transactions[j].NumberofConfirmations = 0
					}
				}
				x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions =
					electrs.GetBalanceCount(elem)
			} else {
				log.Println("error in gettxsaddress call: ", err)
			}
		}
	}

	return x, nil
}

func multiBalance(arr []string, w http.ResponseWriter, r *http.Request) format.BalanceReturn {
	var x format.BalanceReturn
	for _, elem := range arr {
		// send the request out
		tBalance, tUnconfirmedBalance := 0.0, 0.0
		if opts.Mainnet {
			go func(elem string) {
				log.Println("calling the balances endpoint")
				tBalance, tUnconfirmedBalance = electrs.GetBalanceAddress(elem)
				x.Balance += tBalance
				x.UnconfirmedBalance += tUnconfirmedBalance
			}(elem)
			wait()
		} else {
			tBalance, tUnconfirmedBalance = electrs.GetBalanceAddress(elem)
			x.Balance += tBalance
			x.UnconfirmedBalance += tUnconfirmedBalance
		}
	}

	return x
}

// MultiUtxos gets the utxos associated with multiple addresses
func MultiUtxos(router *httprouter.Router) {
	router.POST("/utxos", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var result [][]format.Utxo
		if opts.Mainnet {
			for _, elem := range arr {
				// send the request out
				go func(elem string) {
					tempTxs, err := electrs.GetUtxosAddress(elem)
					if err != nil {
						erpc.ResponseHandler(w, http.StatusInternalServerError)
						log.Println(err)
						return
					}
					result = append(result, tempTxs)
				}(elem)
			}
			wait()
		} else {
			for _, elem := range arr {
				tempTxs, err := electrs.GetUtxosAddress(elem)
				if err != nil {
					erpc.ResponseHandler(w, http.StatusInternalServerError)
					log.Println(err)
					return
				}
				result = append(result, tempTxs)
			}
		}
		erpc.MarshalSend(w, result)
	})
}

// MultiData gets all data associated with a particular address
func MultiData(router *httprouter.Router) {
	router.POST("/data", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// MultiBalTxs combines the balance and multiaddr endpoints
func MultiBalTxs(router *httprouter.Router) {
	router.POST("/baltxs", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var ret format.BalTxReturn

		if opts.Mainnet {
			go func(arr []string) {
				ret.Balance = multiBalance(arr, w, r)
			}(arr)
		} else {
			ret.Balance = multiBalance(arr, w, r)
		}

		// multiAddr is a synch call, so multiBalance should finish before
		ret.Transactions, err = multiAddr(w, r, arr)
		if err != nil {
			return
		}

		erpc.MarshalSend(w, ret)
	})
}

// MultiBalances gets the net balance associated with multiple addresses
func MultiBalances(router *httprouter.Router) {
	router.POST("/balances", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		x := multiBalance(arr, w, r)
		erpc.MarshalSend(w, x)
	})
}

// MultiTxs gets the transactions associated with multiple addresses
func MultiTxs(router *httprouter.Router) {
	router.POST("/txs", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var x format.TxReturn
		for _, elem := range arr {
			// send the request out
			tempTxs, err := electrs.GetTxsAddress(elem)
			if err != nil {
				erpc.ResponseHandler(w, http.StatusInternalServerError)
				log.Println(err)
				return
			}
			x.Txs = append(x.Txs, tempTxs)
		}

		erpc.MarshalSend(w, x)
	})
}

// GetFees gets the current fee estimate from electrs
func GetFees(router *httprouter.Router) {
	router.POST("/fees", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			log.Println(err)
			return
		}

		x, err := electrs.GetFeeEstimates()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			log.Println(err)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// PostTx posts a transaction to the blockchain
func PostTx(router *httprouter.Router) {
	router.POST("/tx", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		body := electrs.ElectrsURL + "/tx"
		data, err := erpc.PostRequest(body, r.Body)
		if err != nil {
			log.Println("could not submit transacation to testnet, quitting")
		}
		var x interface{}
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println("error while unmarshalling json struct", string(data))
			w.Write(data)
			return
		}
		erpc.MarshalSend(w, x)
	})
}
