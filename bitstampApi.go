package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Api interface {
	GetFee() float64
	GetName() string
	GetLast() float64
	RenewLast() float64

	setBitcoin(amount float64) float64
	setFunds(amount float64) float64
}

func (api *BitstampApi) setBitcoin(amount float64) float64 {
	api.bitcoin += amount
	return api.bitcoin
}

func (api *BitstampApi) setFunds(amount float64) float64 {
	api.funds += amount
	return api.funds
}

func newBitstampApi() *BitstampApi {

	api := &BitstampApi{}
	api.client = &http.Client{}

	api.pair = "btceur" //TODO: Hardcoded for now, should be moved to config file
	api.publicURL = "https://www.bitstamp.net/api/v2/"

	api.funds = 300
	api.bitcoin = 0.03

	return api
}

type BitstampApi struct {
	publicURL string
	pair      string
	key       string
	secret    string
	client    *http.Client
	last      float64
	bitcoin   float64
	funds     float64
}

type BitstampTicket struct {
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	Bid       string `json:"bid"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}

//Create request
func (api *BitstampApi) doRequest(parameters map[string]string, url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request failed")
		log.Print(err)
	}
	req.Header.Add("User-Agent", "GoTrade")
	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()
	// Save a copy of this request for debugging.
	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(requestDump))
	resp, err := api.client.Do(req)
	if err != nil {
		fmt.Println("Request failed")
		log.Print(err)
	}

	// Read request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}

	return body
}

func (api *BitstampApi) RenewLast() float64 {
	resp := api.doRequest(nil, api.publicURL+"ticker/btceur")
	ticket := BitstampTicket{}
	json.Unmarshal(resp, &ticket)
	api.last = StringtoFloat(ticket.Last)
	return StringtoFloat(ticket.Last)
}

func (api *BitstampApi) GetLast() float64 {
	return api.last
}

func (api *BitstampApi) GetName() string {
	return "Bitstamp"
}

func (api *BitstampApi) GetFee() float64 {
	return 0.25
}
