package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRequest(url string) ([]byte, error) {
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
	MarshalSend(w, response)
}

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func MarshalSend(w http.ResponseWriter, x interface{}) {
	xJson, err := json.Marshal(x)
	if err != nil {
		log.Println("did not marshal json", err)
		responseHandler(w, http.StatusInternalServerError)
		return
	}
	WriteToHandler(w, xJson)
}

// GetAndSendJsonBalance is a handler that makes a get request and returns json data
func GetBalanceAddress(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := "http://34.73.144.32:443/address/" + addr
	data, err := GetRequest(body)
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

func multigetBalance() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/multigetbalance", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		checkGetRequest(w, r) // check origin of request as well if needed
		var arr []string
		err := r.ParseForm()
		if err != nil {
			log.Println("did not parse form", err)
			return prepProject, err
		}
		addresses := r.FormValue("addresses")
		// this is a comma separated list. Parse that.
		arr := string.Split(addresses, ",")
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
		MarshalSend(w, x)
	})
}

func startHandlers() {
	multigetBalance()
}

func main() {
	startHandlers()
	log.Fatal(http.ListenAndServe("3001", nil)) // 3001 os hardcoded as of now
}
