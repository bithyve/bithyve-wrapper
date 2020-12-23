package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// use this function if you're writing some new endpoints which improve performance
// the new endpoint should be at /new so we can compare benchmarks directly without
// having to run both separately

var APIURL = "https://test-wrapper.bithyve.com"

func postRoutine(url string, inputx string) {
	input := []byte(inputx)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(input))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	log.Println("response Body:", len(string(body)))
}

func BenchmarkMultiAddr(b *testing.B) {
	url := APIURL + "/multiaddr"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkMultiAddrNew(b *testing.B) {
	url := APIURL + "/multiaddrnew"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkBalTx(b *testing.B) {
	url := APIURL + "/baltxs"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkBalTxNew(b *testing.B) {
	url := APIURL + "/baltxsnew"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkMultiGetUtxos(b *testing.B) {
	url := APIURL + "/multigetutxos"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkMultiGetUtxosNew(b *testing.B) {
	url := APIURL + "/multigetutxosnew"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkMultiGetBalance(b *testing.B) {
	url := APIURL + "/multigetbalancenew"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkMultiGetBalanceNew(b *testing.B) {
	url := APIURL + "/multigetbalancenew"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkMultiGetTxs(b *testing.B) {
	url := APIURL + "/multigettxs"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkMultiGetTxsNew(b *testing.B) {
	url := APIURL + "/multigettxsnew"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 20; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}
