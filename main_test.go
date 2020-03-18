package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/bithyve/bithyve-wrapper/format"

	erpc "github.com/Varunram/essentials/rpc"
	"github.com/bithyve/bithyve-wrapper/electrs"
)

var arr = []byte(`{"addresses":["3QTNMWyYJaPTBXT5QHMeH6QMoyPgsesaXs"]}`)
var mainU = "https://api.bithyve.com"
var fallback = "https://blockstream.info/api"

func testUtxos(wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()
	body := mainU + "/utxos"

	data, err := erpc.PostRequest(body, bytes.NewBuffer(arr))
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}

	var x [][]format.Utxo
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}
	log.Println("/utxos endpoint works")
}

func testData(wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()
	body := mainU + "/data"

	data, err := erpc.PostRequest(body, bytes.NewBuffer(arr))
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}

	var x []format.MultigetAddrReturn
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}
	log.Println("/data endpoint works")
}

func testBalTxs(wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()
	body := mainU + "/baltxs"

	data, err := erpc.PostRequest(body, bytes.NewBuffer(arr))
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}

	var x format.BalTxReturn
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}
	log.Println("/baltxs endpoint works")
}

func testBalances(wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()
	body := mainU + "/balances"

	data, err := erpc.PostRequest(body, bytes.NewBuffer(arr))
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}

	var x format.BalanceReturn
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}
	log.Println("/balances endpoint works")
}

func testTxs(wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()
	body := mainU + "/txs"

	data, err := erpc.PostRequest(body, bytes.NewBuffer(arr))
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}

	var x format.TxReturn
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}
	log.Println("/txs endpoint works")
}

func testFees(wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()
	body := mainU + "/fees"
	formdata := url.Values{}
	reader := strings.NewReader(formdata.Encode())

	data, err := erpc.PostRequest(body, reader)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}

	var x format.FeeResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		t.Fatal(err)
	}
	log.Println("/fees endpoint works")
}

func TestMain(t *testing.T) {
	electrs.SetURL(mainU, fallback)
	erpc.SetConsts(10)

	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go testUtxos(&wg, t)
		wg.Add(1)
		go testData(&wg, t)
		wg.Add(1)
		go testBalTxs(&wg, t)
		wg.Add(1)
		go testBalances(&wg, t)
		wg.Add(1)
		go testTxs(&wg, t)
		wg.Add(1)
		go testFees(&wg, t)
	}

	wg.Wait()
	log.Println("test success")
}
