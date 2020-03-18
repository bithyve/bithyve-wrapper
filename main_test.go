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

var arr1 = []byte(`{"addresses":["3QTNMWyYJaPTBXT5QHMeH6QMoyPgsesaXs"]}`)
var arr5 = []byte(`{"addresses":["2ND8uocTbhVQyyYPysbLR2b7o3EhPkSz9C3","2N9RqBCrQw1x83pHT5dyLK6gsbV7r9u3zeF","2N7yYm8E2bxW4PoyRvWwy6Ha2sUCCPDgtWH","2NF9inTAZweCzWwwvpnLtDM3BM5Q6ELHLRL","2N1R8HBkrnvPDGMrtTtRTFpfbZt9SieyMNg"]}`)
var arr10 = []byte(`{"addresses":["2MuF1szXnbaiZmxg51XDCZBBtRQM3ud2GvD","2MxPeBCFnJWnuQsAuQaMk5e4UJ3JApNHoJ8","2NCJhM17H9AkmwUFTzJEoZgKvSgrVEW83Ax","2N9X7hWZ8U4zNmoqXGTm8khwL3FvTkG7F3r","2N8x8FJnN6BsMf9xYtDsw1EPFyZ25QkRVBS","2MvCsp6PVRYmp96MQAvL2x7awUWFqyQnFAd","2N6zahXGBnceNojsRKtJJLW51Q2rf5sCn3k","2MyUVcFieKAHfx9QvTgMwNGeoZhQGvKTWb1","2NDgdTxJJkjBG9DFkjh4Rag8YzaMfJESLMy","2N29Yc7b8wNpT1GGzi72k7TVNBtfzmL43Z7"]}`)
var arr20 = []byte(`{"addresses":["2N716XidW2TF5ME7A91mJHsivWnQZDwCkiK","2Mwu5Vnr2E2vJHpqrXpnAcc23nkpTvEHa8b","2N7WZBptfECLrA9EaveJuXtbyc3hTW4hG83","2Mu9heRDb9kGPTF9eejYNfig5eCEZ6NgXnm","2N8W7Xe2Q5AFBtKnB5pMEoFeGbri8j3Mq8s","2N8UYChaQUUDZes78Ba8FT7zLzRkcdfTJfC","2NCkQqxhokxaKT7PEbThSm2TRAHsnL1ihN5","2MsPvt2z61iuhvT8viEf6xvb47dQQCfDhhc","2N6tqnUDDcZHSX4UJe951EW5ac8eFvpPxM3","2Mznk4P5HjqmW9tTnKogwkkrZn3W5djDCeu","2MzLadTVaWL7BPi5D3ZLWaFUEeatHAQT1nb","2Mwp2W8nCbZokFcJbnYLoSamF1XUqLJuNDm","2N3u787aNsG2GktPCE7CpcRkVGLJztHa4mr","2NGStTcmxuExP2FeMAWoiVXSQ5N8jC3Fnwb","2N4ughHBT8AoBkuLkoXGMVQ9z8tvDQ3JaiZ", "2N7pG1z1GHWLXyZ1o9cKnUyV6oDym9cJBgP","2MwvmGTheh9Mc8kfj7upd6qHD6sm65sXgxT","2N4TdkKSsQorh6r6wL4ezznwEuPUnvpkw8H","2MuBixYRMfF6knrhxr6bNehNjdVmnktdAUk","2MtNy6bcaBzCRNxRv1mYA1RtXE92HtHMff2","2NB1UPUCcbYMc9aM1ZUgsdjKPnbu9oRcbZW","2MuYu5FCbWbpGMcJMSeKtZSv6W7KE4BLaxX","2N22umYZrQnkYDYWbcvVzfv4eS4anK7gTWs","2N8jjguZayhuLqhuQqT8sEHT6in4HniUobn","2N1oRMkDVv5cBaugMG8tXyGKPqe6dJhJ1zY"]}`)
var arr40 = []byte(`{"addresses":["35hTZ8vCYSUCn9E7CPMqTXZhgnRQaVKzQR","38WZB6kwFhJ8ND49uXPgzn58iHbgZmPaUB","38Av9g2HBEt7NztDUhCViTVSm2H8aQKpwX","3ELvizGxn5wiBccqmenrDjsuJm6iUHXZFS","3NzLbZT3Lg1f8ZP7S1GSLHZ763qvyGaty3","35efxt8CJoCBpTqKL1oaeWUkLithZb8Sh5","3Eyz2egRvftWp4vctSvr4LhCh8nyzs1V7H","33fKeXaVx7kzpNfyikGK239PvBg3vinKFP","3N4gje1dog37qvu6do1R4TNydhP2waJjbt","331qd6DdgcpirU74RtGaeqz3cWpUQq7vtc","38pCzDUJYZXyCb5rHdqUJ9Po6xkEheEoHL","3CUMDa6JBxob9Ur1axdoTkLmZcDCDWGuyE","3E2sEgyx2mvTiZ6vprRG8ceQr9AZE8VL8v","3F4TGFnfYvRe7P9GZNWG9QJjCFaQd8hjmd","3FbohTkJUMGftQ15hcxcKM23sWuPwteypU","339tamnwnXN9stg9w4RYor81PgAevMzkEN","327AZm7PBNEL56EXp6PXKTmZBF6PKCku6H","3DET8Jafai5btDdmsd2uFMM3JE9EzprSsr","3FpnAhBQfRdEn63N2yfR8guca9Fuhoyuf4","31ieAUGHHWqVT6F5Z7UvaSyo56JRnxdM7L","3Ksd7hK9sWbnuUJM5zbBDHTTAviE8bnPdN","37Ka8w9hXv99QunvmGuVqwzPZknNqT8x7g","3JrZ17frqAGeQ35R5hnAxN25RKrQZizNoF","393LFaMU9c5M2yvXGFbhDm2TZNmGfZiMDF","3JyAgfzJjtxLVrbrK8qaXZzi8mM6bdxvxc","3AvJ14uq2dXzAAh4pKfpkdTNAn9ehXjMgG","35sLAG3W6EgDepQv8kSkg8DrYAWJVrWCKS","3CZkhR5WaVxuXwRRwfrpUg1K5hzbNVWgCc","3CSS4yZVUPkjaShHK1S25pNy6UYGbcsJKC","31sNSDPQqN5hNJ4cS97JRNKg9PNDq5x1Uv","3E9Pg1WsgKbgSLbjWtSq3K6kPrNRSaQzeS","3DdeYzunFQJA3Fq8nL6ET7Pc1nti1nM792","3GjhyrmKT8xB7aaSQ1hLKifGBksn4X5ncM","37MXbRuNjPBB5cmUhhQjgUTc86JeZkMZnt","35pCjoinyNtiArG4BJBeH7th13pU1tbsXd","3ERUK726qRZnTE3hQCRYaPnjvWotRcpcZy","3A8cu4bpo4dFRkcoipgqawBVDm6vthJBBe","3K3yd3jCuc2TXDp6YJMPvAB4ovxxP5PMQB","33tSJjjoDnDmJGJuXBgSsovnw3SReW3NKn","36rqsrXWueytcg7hfjLHNJNNSBHNZsudCD"]}`)

var mainU = "https://api.bithyve.com"
var fallback = "https://blockstream.info/api"

func testUtxos(wg *sync.WaitGroup, t *testing.T, arr []byte) {
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
	log.Println("/utxos endpoint works: ", len(arr))
}

func testData(wg *sync.WaitGroup, t *testing.T, arr []byte) {
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
	log.Println("/data endpoint works: ", len(arr))
}

func testBalTxs(wg *sync.WaitGroup, t *testing.T, arr []byte) {
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
	log.Println("/baltxs endpoint works: ", len(arr))
}

func testBalances(wg *sync.WaitGroup, t *testing.T, arr []byte) {
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
	log.Println("/balances endpoint works: ", len(arr))
}

func testTxs(wg *sync.WaitGroup, t *testing.T, arr []byte) {
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
	log.Println("/txs endpoint works: ", len(arr))
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

func TestOne(t *testing.T) {
	electrs.SetURL(mainU, fallback)
	erpc.SetConsts(10)

	inputArray := make([][]byte, 5)
	inputArray[0] = arr1
	inputArray[1] = arr5
	inputArray[2] = arr10
	inputArray[3] = arr20
	inputArray[4] = arr40

	log.Println("len arr:", len(inputArray[0]), len(inputArray[1]), len(inputArray[2]), len(inputArray[3]), len(inputArray[4]))
	var wg sync.WaitGroup

	wg.Add(1)
	go testFees(&wg, t)

	for _, arr := range inputArray {
		wg.Add(1)
		go testUtxos(&wg, t, arr)
		wg.Add(1)
		go testData(&wg, t, arr)
		wg.Add(1)
		go testBalTxs(&wg, t, arr)
		wg.Add(1)
		go testBalances(&wg, t, arr)
		wg.Add(1)
		go testTxs(&wg, t, arr)
	}

	wg.Wait()
	log.Println("test success")
}
