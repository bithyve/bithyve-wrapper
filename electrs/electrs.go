package electrs

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bithyve/bithyve-wrapper/format"

	erpc "github.com/Varunram/essentials/rpc"
	"github.com/Varunram/essentials/utils"
)

// ElectrsURL is the URL of a running electrs instance
var ElectrsURL = "http://testapi.bithyve.com"

// CurrentBlockHeight gets the current block height from the blockchain
func CurrentBlockHeight() (float64, error) {
	body := ElectrsURL + "/blocks/tip/height"
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return -1, err
	}

	return utils.ToFloat(data)
}

// GetBalanceCount gets the total incoming balance
func GetBalanceCount(w http.ResponseWriter, r *http.Request, addr string) (float64, float64) {
	body := ElectrsURL + "/address/" + addr
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return -1, -1
	}
	// now data is in byte, we need the other structure now
	var x format.Balance
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
	var x format.Balance
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		return -1, -1
	}

	return x.ChainStats.FundedTxoSum - x.ChainStats.SpentTxoSum,
		x.MempoolStats.FundedTxoSum - x.MempoolStats.SpentTxoSum
}

// GetTxsAddress gets the transactions associated with a given address
func GetTxsAddress(w http.ResponseWriter, r *http.Request, addr string) ([]format.Tx, error) {
	var x []format.Tx
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

// GetUtxosAddress gets the utxos associated with a given address
func GetUtxosAddress(w http.ResponseWriter, r *http.Request, addr string) ([]format.Utxo, error) {
	var x []format.Utxo
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
