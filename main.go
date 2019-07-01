package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	//"strings"

	erpc "github.com/Varunram/essentials/rpc"
)

func GetBalanceCount(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := "http://testapi.bithyve.com/address/" + addr
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return -1, -1
	}
	// now data is in byte, we need the other structure now
	var x GetBalanceFormat
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		return -1, -1
	}

	return x.ChainStats.Funded_txo_count, x.MempoolStats.Funded_txo_count
}

func GetBalanceAddress(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := "http://testapi.bithyve.com/address/" + addr
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return -1, -1
	}
	// now data is in byte, we need the other structure now
	var x GetBalanceFormat
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		return -1, -1
	}

	return x.ChainStats.Funded_txo_sum, x.MempoolStats.Funded_txo_sum
}

type Tx struct {
	Txid     string  `json:"txid"`
	Version  float64 `json:"version"`
	Locktime float64 `json:"locktime"`
	Vin      []struct {
		Txid    string  `json:"txid"`
		Vout    float64 `json:"vout"`
		PrevOut struct {
			Scriptpubkey         string  `json:"scriptpubkey"`
			Scriptpubkey_asm     string  `json:"scriptpubkey_asm"`
			Scriptpubkey_address string  `json:"scriptpubkey_address"`
			Scriptpubkey_type    string  `json:"scriptpubkey_type"`
			Value                float64 `json:"value"`
		} `json:"prevout"`
		Scriptsig     string   `json:"scriptsig"`
		Scriptsig_asm string   `json:"scriptsig_asm"`
		Witness       []string `json:"witness"`
		Is_coinbase   bool     `json:"is_coinbase"`
		Sequence      float64  `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Scriptpubkey         string  `json:"scriptpubkey"`
		Scriptpubkey_asm     string  `json:"scriptpubkey_asm"`
		Scriptpubkey_address string  `json:"scriptpubkey_address"`
		Scriptpubkey_type    string  `json:"scriptpubkey_type"`
		Value                float64 `json:"value"`
	}
	Size   float64 `json:"size"`
	weight float64 `json:"weight"`
	Fee    float64 `json:"fee"`
	Status struct {
		Confirmed    bool    `json:"confirmed"`
		Block_height float64 `json:"block_height"`
		Block_hash   string  `json:"block_hash"`
		Block_time   float64 `json:"block_time"`
	}
	NumberofConfirmations float64
}

func GetTxsAddress(w http.ResponseWriter, r *http.Request, addr string) ([]Tx, error) {
	var x []Tx
	body := "http://testapi.bithyve.com/address/" + addr + "/txs"
	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return x, err
	}
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		return x, err
	}

	return x, nil
}

type UtxoVout struct {
	Scriptpubkey         string  `json:"scriptpubkey"`
	Scriptpubkey_asm     string  `json:"scriptpubkey_asm"`
	Scriptpubkey_address string  `json:"scriptpubkey_address"`
	Scriptpubkey_type    string  `json:"scriptpubkey_type"`
	Value                float64 `json:"value"`
	Index                int
	Address              string
}

type Utxo struct {
	Txid   string `json:"txid"`
	Vout   int    `json:"vout"`
	Status struct {
		Confirmed    bool    `json:"confirmed"`
		Block_height float64 `json:"block_height"`
		Block_hash   string  `json:"block_hash"`
		Block_time   float64 `json:"block_time"`
	} `json:"status"`
	Value   float64 `json:"value"`
	Address string
}

func GetUtxosAddress(w http.ResponseWriter, r *http.Request, addr string) ([]Utxo, error) {
	var x []Utxo
	body := "http://testapi.bithyve.com/address/" + addr + "/utxo"
	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return nil, err
	}
	// now data is in byte, we need the other structure now
	log.Println(string(data))
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		return nil, err
	}

	for i, _ := range x {
		x[i].Address = addr
	}
	return x, nil
}

type GetBalanceFormat struct {
	Address    string `json:"address"`
	ChainStats struct {
		Funded_txo_count float64 `json:"funded_txo_count"`
		Funded_txo_sum   float64 `json:"funded_txo_sum"`
		Spent_txo_count  float64 `json:"spent_txo_count"`
		Spent_txo_sum    float64 `json:"spent_txo_sum"`
		Tx_count         float64 `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		Funded_txo_count float64 `json:"funded_txo_count"`
		Funded_txo_sum   float64 `json:"funded_txo_sum"`
		Spent_txo_count  float64 `json:"spent_txo_count"`
		Spent_txo_sum    float64 `json:"spent_txo_sum `
		Tx_count         float64 `json:"tx_count `
	} `json:"mempool_stats"`
}

type MultigetBalanceReturn struct {
	Balance            float64
	UnconfirmedBalance float64
}

type RequestFormat struct {
	Addresses []string `json:"addresses"`
}

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
		balance := float64(0)
		uBalance := float64(0)
		for _, elem := range arr {
			// send the request out
			tBalance, tUnconfirmedBalance := GetBalanceAddress(w, r, elem)
			balance += tBalance
			uBalance += tUnconfirmedBalance
		}
		var x MultigetBalanceReturn
		x.Balance = balance
		x.UnconfirmedBalance = uBalance
		erpc.MarshalSend(w, x)
	})
}

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
		var result [][]Tx
		for _, elem := range arr {
			// send the request out
			tempTxs, err := GetTxsAddress(w, r, elem)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, http.StatusInternalServerError)
				return
			}
			result = append(result, tempTxs)
		}
		erpc.MarshalSend(w, result)
	})
}

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
			tempTxs, err := GetUtxosAddress(w, r, elem)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, http.StatusInternalServerError)
				return
			}
			result = append(result, tempTxs)
		}
		erpc.MarshalSend(w, result)
	})
}

type MultigetAddrReturn struct {
	TotalTransactions       float64
	ConfirmedTransactions   float64
	UnconfirmedTransactions float64
	Transactions            []Tx
	Address                 string
}

func CurrentBlockHeight() (float64, error) {
	body := "http://testapi.bithyve.com/blocks/tip/height"
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return -1, err
	}

	// now the data needs to be converted into an integer ie string to float
	stringBn := string(data)
	intBn, err := strconv.ParseFloat(stringBn, 32)
	if err != nil {
		return -1, err
	}
	return intBn, nil
}

func MultigetAddr() {
	// make a curl request out to lcoalhost and get the ping response
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
			allTxs, err := GetTxsAddress(w, r, elem)
			if err != nil {
				continue
			}
			x[i].TotalTransactions = float64(len(allTxs))
			x[i].Transactions = allTxs
			for j, _ := range x[i].Transactions {
				x[i].Transactions[j].NumberofConfirmations = currentBh - x[i].Transactions[j].Status.Block_height
			}
			x[i].ConfirmedTransactions, x[i].UnconfirmedTransactions = GetBalanceCount(w, r, elem)
		}
		erpc.MarshalSend(w, x)
	})
}

type FeeResponse struct {
	Two              float64 `json:"2"`
	Three            float64 `json:"3"`
	Four             float64 `json:"4"`
	Six              float64 `json:"6"`
	Ten              float64 `json:"10"`
	Twenty           float64 `json:"20"`
	OneFourFour      float64 `json:"144"`
	FiveZeroFour     float64 `json:"504"`
	OneThousandEight float64 `json:"1008"`
}

func GetFees() {
	http.HandleFunc("/fees", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}

		var x FeeResponse
		body := "http://testapi.bithyve.com/fee-estimates"
		erpc.GetAndSendJson(w, body, x)
	})
}

func PostTx() {
	http.HandleFunc("/tx", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		body := "http://testapi.bithyve.com/tx"
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

func RelayTxid() {
	http.HandleFunc("/txid", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		if r.URL.Query()["txid"] == nil {
			erpc.ResponseHandler(w, http.StatusBadRequest)
			return
		}

		txid := r.URL.Query()["txid"][0]
		body := "http://testapi.bithyve.com/tx/" + txid
		var x Tx
		erpc.GetAndSendJson(w, body, x)
	})
}

func RelayGetRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckGet(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		// log.Println(r.URL.String())
		body := "http://testapi.bithyve.com" + r.URL.String()
		data, err := erpc.GetRequest(body)
		if err != nil {
			log.Println("could not submit transacation to testnet, quitting")
			erpc.ResponseHandler(w, http.StatusInternalServerError)
			return
		}

		var x interface{}
		_ = json.Unmarshal(data, &x)
		erpc.MarshalSend(w, x)
	})
}

func startHandlers() {
	MultigetBalance()
	MultigetTxs()
	MultigetUtxos()
	MultigetAddr()
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
