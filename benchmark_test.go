package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

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
	url := "https://testapi.bithyve.com/multiaddr"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 10; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}

func BenchmarkBalTxs(b *testing.B) {
	url := "https://testapi.bithyve.com/baltxs"
	input1 := `{"addresses":["2MsxyDNd4kMiRxi8PbXVPvuk526fAWRAaSD", "2N7dRtWLBJgC7QdmEaSLNyiJtfrnvJtanMb"]}`

	log.Println("Called")
	b.StartTimer()
	for i := 0; i < 10; i++ {
		postRoutine(url, input1)
	}
	b.StopTimer()
}
