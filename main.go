package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func Get(url string) ([]byte, error) {
	var dummy []byte
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("did not create new GET request", err)
		return dummy, err
	}
	req.Header.Set("Origin", "localhost")
	res, err := client.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

type StatusResponse struct {
	Code   int
	Status string
}

func responseHandler(w http.ResponseWriter, status int) {
	var response StatusResponse
	response.Code = status
	switch status {
	case http.StatusOK:
		response.Status = "OK"
	case http.StatusBadRequest:
		response.Status = "Bad Request error!"
	case http.StatusNotFound:
		response.Status = "404 Error Not Found!"
	case http.StatusInternalServerError:
		response.Status = "Internal Server Error"
	default:
		response.Status = "404 Page Not Found"
	}
	Send(w, response)
}

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func Send(w http.ResponseWriter, x interface{}) {
	xJson, err := json.Marshal(x)
	if err != nil {
		log.Println("did not marshal json", err)
		responseHandler(w, http.StatusInternalServerError)
		return
	}
	WriteToHandler(w, xJson)
}

type Status404 struct {
	Status int
}

func Send404(w http.ResponseWriter) {
	var x Status404
	x.Status = 404
	Send(w, x)
}

// GetAndSendJsonBalance is a handler that makes a get request and returns json data
func GetBalanceCount(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := "http://35.229.68.185:443/address/" + addr
	data, err := Get(body)
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

// GetAndSendJsonBalance is a handler that makes a get request and returns json data
func GetBalanceAddress(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := "http://35.229.68.185:443/address/" + addr
	data, err := Get(body)
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
	body := "http://35.229.68.185:443/address/" + addr + "/txs"
	log.Println(body)
	data, err := Get(body)
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

// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" -d 'addresses=17rdSE552fTwvRqLxdKJtfkncB1om8XtJT%2C17rdSE552fTwvRqLxdKJtfkncB1om8XtJT%2C17rdSE552fTwvRqLxdKJtfkncB1om8XtJT' "http://localhost:3001/multigetutxos"
func GetUtxosAddress(w http.ResponseWriter, r *http.Request, addr string) ([]Utxo, error) {
	var x []Utxo
	body := "http://35.229.68.185:443/address/" + addr + "/utxo"
	log.Println(body)
	data, err := Get(body)
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

func checkGetRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func checkPostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

type MultigetBalance struct {
	Balance            float64
	UnconfirmedBalance float64
}

type RequestFormat struct {
	Addresses []string `json:"addresses"`
}

// example request:
// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" -d 'addresses=17rdSE552fTwvRqLxdKJtfkncB1om8XtJT%2C17rdSE552fTwvRqLxdKJtfkncB1om8XtJT%2C17rdSE552fTwvRqLxdKJtfkncB1om8XtJT' "http://localhost:3001/multigetbalance"
func multigetBalance() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigetbalance", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		checkPostRequest(w, r) // check origin of request as well if needed
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			Send404(w)
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			Send404(w)
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
		var x MultigetBalance
		x.Balance = balance
		x.UnconfirmedBalance = uBalance
		Send(w, x)
	})
}

// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" -d 'addresses=17rdSE552fTwvRqLxdKJtfkncB1om8XtJT%2C17rdSE552fTwvRqLxdKJtfkncB1om8XtJT%2C17rdSE552fTwvRqLxdKJtfkncB1om8XtJT' "http://localhost:3001/multigettxs"
func multigetTxs() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigettxs", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		checkPostRequest(w, r) // check origin of request as well if needed
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			Send404(w)
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			Send404(w)
		}
		arr := rf.Addresses
		var result [][]Tx
		for _, elem := range arr {
			// send the request out
			tempTxs, err := GetTxsAddress(w, r, elem)
			if err != nil {
				log.Println(err)
				responseHandler(w, http.StatusInternalServerError)
				return
			}
			result = append(result, tempTxs)
		}
		Send(w, result)
	})
}

func multigetUtxos() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigetutxos", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		checkPostRequest(w, r) // check origin of request as well if needed
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			Send404(w)
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			Send404(w)
		}
		arr := rf.Addresses
		var result [][]Utxo
		for _, elem := range arr {
			// send the request out
			tempTxs, err := GetUtxosAddress(w, r, elem)
			if err != nil {
				log.Println(err)
				responseHandler(w, http.StatusInternalServerError)
				return
			}
			result = append(result, tempTxs)
		}
		Send(w, result)
	})
}

type MultigetAddr struct {
	TotalTransactions       float64
	ConfirmedTransactions   float64
	UnconfirmedTransactions float64
	Transactions            []Tx
	Address                 string
}

func currentBlockHeight() (float64, error) {
	body := "http://35.229.68.185:443/blocks/tip/height"
	data, err := Get(body)
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

func multigetAddr() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multiaddr", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		checkPostRequest(w, r) // check origin of request as well if needed
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			Send404(w)
			return
		}
		var rf RequestFormat
		err = json.Unmarshal(data, &rf)
		if err != nil {
			log.Println(err)
			Send404(w)
			return
		}
		arr := rf.Addresses
		x := make([]MultigetAddr, len(arr))
		currentBh, err := currentBlockHeight()
		if err != nil {
			log.Println(err)
			Send404(w)
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
		Send(w, x)
	})
}

func startHandlers() {
	multigetBalance()
	multigetTxs()
	multigetUtxos()
	multigetAddr()
}

func main() {
	startHandlers()
	// if you're running esplora, use socat tcp-listen:3003,reuseaddr,fork tcp:localhost:3002 to tunnel port since
	// it does not seem possible to open the port directly
	log.Fatal(http.ListenAndServe("localhost:3002", nil)) // 3001 os hardcoded as of now
}
