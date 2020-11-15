package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bithyve/bithyve-wrapper/electrs"
	"github.com/bithyve/bithyve-wrapper/format"

	erpc "github.com/Varunram/essentials/rpc"
)

var (
	// APIError is the response message returned if there's something wrong with electrs
	APIError = "API error, couldn't contact electrs"
	// JSONError is the response message returned if there's a problem with converting a bytestring to json
	JSONError = "Error while converting response to json"
)

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
		log.Println(err)
		return arr, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		log.Println(err)
		return arr, err
	}
	var rf format.RequestFormat
	err = json.Unmarshal(data, &rf)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError, JSONError)
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

func checkReqEI(w http.ResponseWriter, r *http.Request) (format.EIRequestFormat, error) {
	var earr, iarr []string
	var rf format.EIRequestFormat
	err := erpc.CheckPost(w, r)
	if err != nil {
		log.Println(err)
		return rf, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		log.Println(err)
		return rf, err
	}
	err = json.Unmarshal(data, &rf)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError, JSONError)
		log.Println(err)
		return rf, err
	}

	// now we have an array of ids, earr, iarr. Need to loop through them
	for key, elem := range rf {
		earr = elem.ExternalAddresses
		iarr = elem.InternalAddresses
		// filter through list to remove duplicates
		nodupsMap1 := make(map[string]bool)
		var nodups1 []string

		for _, elem := range earr {
			if _, value := nodupsMap1[elem]; !value {
				nodupsMap1[elem] = true
				nodups1 = append(nodups1, elem)
			}
		}

		// filter through list to remove duplicates
		nodupsMap2 := make(map[string]bool)
		var nodups2 []string

		for _, elem := range iarr {
			if _, value := nodupsMap2[elem]; !value {
				nodupsMap2[elem] = true
				nodups2 = append(nodups2, elem)
			}
		}
		var temp struct {
			ExternalAddresses []string `json:"External"`
			InternalAddresses []string `json:"Internal"`
		}
		temp.ExternalAddresses = nodups1
		temp.InternalAddresses = nodups2
		rf[key] = temp
	}

	return rf, nil
}

func addrHelper(wg *sync.WaitGroup, x []format.MultigetAddrReturn, i int, elem string, currentBh float64) {
	defer wg.Done()

	allTxs, err := electrs.GetTxsAddress(elem)
	if err == nil {
		x[i].TotalTransactions = float64(len(allTxs))
		x[i].Transactions = allTxs
		x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = 0, 0
		for j := range x[i].Transactions {
			if x[i].Transactions[j].Status.Confirmed {
				x[i].Transactions[j].NumberofConfirmations =
					currentBh - x[i].Transactions[j].Status.BlockHeight + 1
			} else {
				x[i].Transactions[j].NumberofConfirmations = 0
			}
		}
	}
}

func addrbalHelper(wg *sync.WaitGroup, x []format.MultigetAddrReturn, i int, elem string) {
	defer wg.Done()
	x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions =
		electrs.GetBalanceCount(elem)
}

func multiAddr(w http.ResponseWriter, r *http.Request,
	arr []string) ([]format.MultigetAddrReturn, error) {

	x := make([]format.MultigetAddrReturn, len(arr))
	currentBh, err := electrs.CurrentBlockHeight()
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError, APIError)
		log.Println(err)
		return x, err
	}

	if opts.Mainnet {
		var wg1 sync.WaitGroup
		var wg2 sync.WaitGroup

		for i, elem := range arr {
			x[i].Address = elem // store the address of the passed elements
			wg1.Add(1)
			go addrHelper(&wg1, x, i, elem, currentBh)
		}

		wg1.Wait()

		for i, elem := range arr {
			wg2.Add(1)
			go addrbalHelper(&wg2, x, i, elem)
		}

		wg2.Wait()
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
							currentBh - x[i].Transactions[j].Status.BlockHeight + 1
					} else {
						x[i].Transactions[j].NumberofConfirmations = 0
					}
					x[i].Transactions[j].Vin = nil
					x[i].Transactions[j].Vout = nil
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

func multiAddrEI(w http.ResponseWriter, r *http.Request,
	earr, iarr []string) ([]format.MultigetAddrReturn, error) {

	var arr []string
	arr = append(earr, iarr...)

	x := make([]format.MultigetAddrReturn, len(arr))
	currentBh, err := electrs.CurrentBlockHeight()
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusInternalServerError, APIError)
		log.Println(err)
		return x, err
	}

	if opts.Mainnet {
		var wg1 sync.WaitGroup
		var wg2 sync.WaitGroup

		for i, elem := range arr {
			x[i].Address = elem // store the address of the passed elements
			wg1.Add(1)
			go addrHelper(&wg1, x, i, elem, currentBh)
		}

		wg1.Wait()

		for i, elem := range arr {
			wg2.Add(1)
			go addrbalHelper(&wg2, x, i, elem)
		}

		wg2.Wait()
	} else {
		var wg4 sync.WaitGroup
		for i, elem := range arr {
			wg4.Add(1)
			go func(wg *sync.WaitGroup, x []format.MultigetAddrReturn, i int, elem string) {
				defer wg.Done()
				allTxs, err := electrs.GetTxsAddress(elem)
				if err == nil {
					if len(allTxs) != 0 {
						x[i].Address = elem
						x[i].TotalTransactions = float64(len(allTxs))
						x[i].Transactions = allTxs
						x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = 0, 0
						for j := range x[i].Transactions {
							if x[i].Transactions[j].Status.Confirmed {
								x[i].Transactions[j].NumberofConfirmations =
									currentBh - x[i].Transactions[j].Status.BlockHeight + 1
							} else {
								x[i].Transactions[j].NumberofConfirmations = 0
							}
						}
						var wg3 sync.WaitGroup
						for j := range x[i].Transactions {
							wg3.Add(1)
							go func(wg *sync.WaitGroup, x []format.MultigetAddrReturn, j int) {
								defer wg.Done()
								x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions =
									electrs.GetBalanceCount(elem)
							}(&wg3, x, j)
						}
						wg3.Wait()
						var wg6 sync.WaitGroup
						for j := range x[i].Transactions {
							wg6.Add(1)
							go func(wg *sync.WaitGroup, x []format.MultigetAddrReturn, j int) {
								defer wg.Done()
								x[i].Transactions[j].Categorize(earr, append(earr, iarr...))
							}(&wg6, x, j)
						}
						wg6.Wait()
					}
				} else {
					log.Println("error in gettxsaddress call: ", err)
				}
			}(&wg4, x, i, elem)
		}
		wg4.Wait()
	}

	var y []format.MultigetAddrReturn

	for _, elem := range x {
		if elem.Address != "" {
			y = append(y, elem)
		}
	}
	return y, nil
}

func balHelper(wg *sync.WaitGroup, elem string, x *format.BalanceReturn) {
	defer wg.Done()
	temp1, temp2 := electrs.GetBalanceAddress(elem)
	x.Balance += temp1
	x.UnconfirmedBalance += temp2
}

func multiBalance(arr []string, w http.ResponseWriter, r *http.Request) format.BalanceReturn {
	if opts.Mainnet {
		var x format.BalanceReturn
		var wg sync.WaitGroup

		for _, elem := range arr {
			wg.Add(1)
			go balHelper(&wg, elem, &x)
		}

		wg.Wait()
		return x
	}

	var x format.BalanceReturn
	for _, elem := range arr {
		tBalance, tUnconfirmedBalance := electrs.GetBalanceAddress(elem)
		x.Balance += tBalance
		x.UnconfirmedBalance += tUnconfirmedBalance
	}
	return x
}

func utxoHelper(wg *sync.WaitGroup, result [][]format.Utxo, i int, elem string) {
	defer wg.Done()
	tempTxs, err := electrs.GetUtxosAddress(elem)
	if err != nil {
		log.Println(err)
		return
	}
	result[i] = tempTxs
}

// MultiUtxos gets the utxos associated with multiple addresses
func MultiUtxos() {
	http.HandleFunc("/utxos", func(w http.ResponseWriter, r *http.Request) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var wg sync.WaitGroup
		result := make([][]format.Utxo, len(arr))
		if opts.Mainnet {
			for i, elem := range arr {
				// send the request out
				wg.Add(1)
				go utxoHelper(&wg, result, i, elem)
			}

			wg.Wait()
			erpc.MarshalSend(w, result)
		} else {
			for i, elem := range arr {
				tempTxs, err := electrs.GetUtxosAddress(elem)
				if err != nil {
					erpc.ResponseHandler(w, http.StatusInternalServerError, APIError)
					log.Println(err)
					return
				}
				result[i] = tempTxs
			}
			erpc.MarshalSend(w, result)
		}
	})
}

// MultiData gets all data associated with a particular address
func MultiData() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
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
func MultiBalTxs() {
	http.HandleFunc("/baltxs", func(w http.ResponseWriter, r *http.Request) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var ret format.BalTxReturn
		ret.Balance = multiBalance(arr, w, r)
		// multiAddr is a synch call, so multiBalance should finish before
		ret.Transactions, err = multiAddr(w, r, arr)
		if err != nil {
			return
		}
		erpc.MarshalSend(w, ret)
	})
}

// MultiUtxoTxs combines the utxo and multiaddr endpoints
func MultiUtxoTxs() {
	http.HandleFunc("/utxotxs", func(w http.ResponseWriter, r *http.Request) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var wg sync.WaitGroup
		var ret format.UtxoTxReturn
		ret.Utxos = make([][]format.Utxo, len(arr))

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			ret.Transactions, err = multiAddr(w, r, arr)
			if err != nil {
				return
			}
		}(&wg)

		for i, elem := range arr {
			// send the request out
			wg.Add(1)
			go utxoHelper(&wg, ret.Utxos, i, elem)
		}

		wg.Wait()
		erpc.MarshalSend(w, ret)
	})
}

// NewMultiUtxoTxs is a new endpoint
func NewMultiUtxoTxs() {
	http.HandleFunc("/nutxotxs", func(w http.ResponseWriter, r *http.Request) {
		rf, err := checkReqEI(w, r)
		if err != nil {
			return
		}

		ret := make(format.EIUtxoReturn, len(rf))
		var wg2 sync.WaitGroup
		for key, elem := range rf {
			wg2.Add(1)
			go func(wg2 *sync.WaitGroup, key string, elem format.EIHelper) {
				defer wg2.Done()
				var arr []string
				iarr := elem.InternalAddresses
				earr := elem.ExternalAddresses
				arr = append(earr, iarr...)
				var wg sync.WaitGroup
				var err error
				var temp format.UtxoTxReturn
				tempUtxos := make([][]format.Utxo, len(earr)+len(iarr))
				// ret.Utxos = make([][]format.Utxo, len(earr)+len(iarr))
				wg.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					temp.Transactions, err = multiAddrEI(w, r, earr, iarr)
					if err != nil {
						return
					}
				}(&wg)

				for i, elem := range arr {
					// send the request out
					wg.Add(1)
					go utxoHelper(&wg, tempUtxos, i, elem)
				}

				wg.Wait()
				// we have both utxos and txs now
				for _, elem := range tempUtxos {
					if len(elem) != 0 {
						temp.Utxos = append(temp.Utxos, elem)
					}
				}
				ret[key] = temp
			}(&wg2, key, elem)
		}
		wg2.Wait()
		erpc.MarshalSend(w, ret)
	})
}

// MultiBalances gets the net balance associated with multiple addresses
func MultiBalances() {
	http.HandleFunc("/balances", func(w http.ResponseWriter, r *http.Request) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		x := multiBalance(arr, w, r)
		erpc.MarshalSend(w, x)
	})
}

func txsHelper(wg *sync.WaitGroup, x *format.TxReturn, i int, elem string) {
	defer wg.Done()
	tempTxs, err := electrs.GetTxsAddress(elem)
	if err != nil {
		log.Println(err)
		return
	}
	x.Txs[i] = make([]format.Tx, len(tempTxs))
	x.Txs[i] = tempTxs
}

// MultiTxs gets the transactions associated with multiple addresses
func MultiTxs() {
	http.HandleFunc("/txs", func(w http.ResponseWriter, r *http.Request) {
		arr, err := checkReq(w, r)
		if err != nil {
			return
		}

		var x format.TxReturn
		var wg sync.WaitGroup
		x.Txs = make([][]format.Tx, len(arr))

		for i, elem := range arr {
			// send the request out
			wg.Add(1)
			go txsHelper(&wg, &x, i, elem)
		}

		wg.Wait()
		erpc.MarshalSend(w, x)
	})
}

// GetFees gets the current fee estimate from electrs
func GetFees(mainnet bool) {
	http.HandleFunc("/fees", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if !mainnet {
			var temp format.FeeResponse
			temp.Two = 5.0
			temp.Three = 4.2
			temp.Four = 3.9
			temp.Five = 3.1
			temp.Six = 2.8
			temp.Ten = 2.0
			temp.Twenty = 1.7
			temp.TwentyFive = 1.6
			temp.OneFourFour = 1.1
			temp.FiveZeroFour = 1.01
			temp.OneThousandEight = 1.0

			erpc.MarshalSend(w, temp)
			return
		}

		x, err := electrs.GetFeeEstimates()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError, APIError)
			log.Println(err)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// GetFeesE gets the current fee estimate from electrs
func GetFeesE(mainnet bool) {
	http.HandleFunc("/fee-estimates", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if !mainnet {
			var temp format.FeeResponse
			temp.Two = 5.0
			temp.Three = 4.2
			temp.Four = 3.9
			temp.Five = 3.1
			temp.Six = 2.8
			temp.Ten = 2.0
			temp.Twenty = 1.7
			temp.TwentyFive = 1.6
			temp.OneFourFour = 1.1
			temp.FiveZeroFour = 1.01
			temp.OneThousandEight = 1.0

			erpc.MarshalSend(w, temp)
			return
		}

		x, err := electrs.GetFeeEstimates()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError, APIError)
			log.Println(err)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// PostTx posts a transaction to the blockchain
func PostTx() {
	http.HandleFunc("/tx", func(w http.ResponseWriter, r *http.Request) {
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
