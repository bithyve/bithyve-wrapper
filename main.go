package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	//"strings"

	erpc "github.com/Varunram/essentials/rpc"
)

// ElectrsURL is the URL of a running electrs instance
var ElectrsURL = "http://testapi.bithyve.com"

// GetBalanceCount gets the total incoming balance
func GetBalanceCount(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := ElectrsURL + "/address/" + addr
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

	return x.ChainStats.FundedTxoCount, x.MempoolStats.FundedTxoCount
}

// GetBalanceAddress gets the net balance of an address
func GetBalanceAddress(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := ElectrsURL + "/address/" + addr
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

	return x.ChainStats.FundedTxoSum - x.ChainStats.SpentTxoSum, x.MempoolStats.FundedTxoSum - x.MempoolStats.SpentTxoSum
}

// Tx is a copy of the transaction structure used by esplora
type Tx struct {
	Txid     string  `json:"txid"`
	Version  float64 `json:"version"`
	Locktime float64 `json:"locktime"`
	Vin      []struct {
		Txid    string  `json:"txid"`
		Vout    float64 `json:"vout"`
		PrevOut struct {
			Scriptpubkey        string  `json:"scriptpubkey"`
			ScriptpubkeyAsm     string  `json:"scriptpubkey_asm"`
			ScriptpubkeyAddress string  `json:"scriptpubkey_address"`
			ScriptpubkeyType    string  `json:"scriptpubkey_type"`
			Value               float64 `json:"value"`
		} `json:"prevout"`
		Scriptsig    string   `json:"scriptsig"`
		ScriptsigAsm string   `json:"scriptsig_asm"`
		Witness      []string `json:"witness"`
		IsCoinbase   bool     `json:"is_coinbase"`
		Sequence     float64  `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Scriptpubkey        string  `json:"scriptpubkey"`
		ScriptpubkeyAsm     string  `json:"scriptpubkey_asm"`
		ScriptpubkeyAddress string  `json:"scriptpubkey_address"`
		ScriptpubkeyType    string  `json:"scriptpubkey_type"`
		Value               float64 `json:"value"`
	}
	Size   float64 `json:"size"`
	Weight float64 `json:"weight"`
	Fee    float64 `json:"fee"`
	Status struct {
		Confirmed   bool    `json:"confirmed"`
		BlockHeight float64 `json:"block_height"`
		BlockHash   string  `json:"block_hash"`
		BlockTime   float64 `json:"block_time"`
	}
	NumberofConfirmations float64
}

// GetTxsAddress gets the transactions associated with a given address
func GetTxsAddress(w http.ResponseWriter, r *http.Request, addr string) ([]Tx, error) {
	var x []Tx
	body := ElectrsURL + "/address/" + addr + "/txs"
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

// UtxoVout is a structure for output utxos
type UtxoVout struct {
	Scriptpubkey        string  `json:"scriptpubkey"`
	ScriptpubkeyAsm     string  `json:"scriptpubkey_asm"`
	ScriptpubkeyAddress string  `json:"scriptpubkey_address"`
	ScriptpubkeyType    string  `json:"scriptpubkey_type"`
	Value               float64 `json:"value"`
	Index               int
	Address             string
}

// Utxo is a copy of the esplora utxo struct
type Utxo struct {
	Txid   string `json:"txid"`
	Vout   int    `json:"vout"`
	Status struct {
		Confirmed   bool    `json:"confirmed"`
		BlockHeight float64 `json:"block_height"`
		BlockHash   string  `json:"block_hash"`
		BlockTime   float64 `json:"block_time"`
	} `json:"status"`
	Value   float64 `json:"value"`
	Address string
}

// GetUtxosAddress gets the utxos associated with a given address
func GetUtxosAddress(w http.ResponseWriter, r *http.Request, addr string) ([]Utxo, error) {
	var x []Utxo
	body := ElectrsURL + "/address/" + addr + "/utxo"
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

	for i := range x {
		x[i].Address = addr
	}
	return x, nil
}

// GetBalanceFormat is a struct that us used to get the blanace from esplora
type GetBalanceFormat struct {
	Address    string `json:"address"`
	ChainStats struct {
		FundedTxoCount float64 `json:"funded_txo_count"`
		FundedTxoSum   float64 `json:"funded_txo_sum"`
		SpentTxoCount  float64 `json:"spent_txo_count"`
		SpentTxoSum    float64 `json:"spent_txo_sum"`
		TxCount        float64 `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		FundedTxoCount float64 `json:"funded_txo_count"`
		FundedTxoSum   float64 `json:"funded_txo_sum"`
		SpentTxoCount  float64 `json:"spent_txo_count"`
		SpentTxoSum    float64 `json:"spent_txo_sum"`
		TxCount        float64 `json:"tx_count"`
	} `json:"mempool_stats"`
}

// MultigetBalanceReturn is a structure that is used for getting multiple balances
type MultigetBalanceReturn struct {
	Balance            float64
	UnconfirmedBalance float64
}

// RequestFormat is the format in which incoming requests hsould arrive for the wrapper to process
type RequestFormat struct {
	Addresses []string `json:"addresses"`
}

// TxReturn is used to return Txs
type TxReturn struct {
	Txs [][]Tx `json:"Txs"`
}

// MultigetAddrReturn is a structure used for multiple addresses json return
type MultigetAddrReturn struct {
	TotalTransactions       float64
	ConfirmedTransactions   float64
	UnconfirmedTransactions float64
	Transactions            []Tx
	Address                 string
}

// CurrentBlockHeight gets the current block height from the blockchain
func CurrentBlockHeight() (float64, error) {
	body := ElectrsURL + "/blocks/tip/height"
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

// FeeResponse is a struct that is returned when a fee query is made
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

// GetFees gets the current fee estimate from esplora
func GetFees() {
	http.HandleFunc("/fees", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckPost(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}

		var x FeeResponse
		body := ElectrsURL + "/fee-estimates"
		erpc.GetAndSendJson(w, body, x)
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
		body := ElectrsURL + "/tx"
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
			erpc.ResponseHandler(w, http.StatusBadRequest)
			return
		}

		txid := r.URL.Query()["txid"][0]
		body := ElectrsURL + "/tx/" + txid
		var x Tx
		erpc.GetAndSendJson(w, body, x)
	})
}

// RelayGetRequest relays all remaining get requests to the esplora instance
func RelayGetRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckGet(w, r) // check origin of request as well if needed
		if err != nil {
			log.Println(err)
			return
		}
		// log.Println(r.URL.String())
		body := ElectrsURL + "" + r.URL.String()
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

// BalTx is a struct used for the baltxs endpoint
type BalTx struct {
	Balance      MultigetBalanceReturn
	Transactions []MultigetAddrReturn `json:"Txs"`
}

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
