package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func newCexioApi() *CexioApi {

	api := &CexioApi{}
	api.client = &http.Client{}

	api.pair = "btceur" //TODO: Hardcoded for now, should be moved to config file
	api.publicURL = "https://cex.io/api/"

	api.funds = 300
	api.bitcoin = 0.03

	return api
}

type CexioApi struct {
	publicURL string
	pair      string
	key       string
	secret    string
	client    *http.Client
	last      float64
	bitcoin   float64
	funds     float64
}

type CexioTicket struct {
	Curr1  string `json:"curr1"`
	Curr2  string `json:"curr2"`
	Lprice string `json:"lprice"`
}

//Create request
func (api *CexioApi) doRequest(parameters map[string]string, url string) []byte {
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

func (api *CexioApi) RenewLast() float64 {
	resp := api.doRequest(nil, api.publicURL+"last_price/BTC/EUR")
	ticket := CexioTicket{}
	json.Unmarshal(resp, &ticket)
	api.last = StringtoFloat(ticket.Lprice)
	return StringtoFloat(ticket.Lprice)
}

func (api *CexioApi) GetLast() float64 {
	return api.last
}

func (api *CexioApi) GetName() string {
	return "Cexio"
}

func (api *CexioApi) GetFee() float64 {
	return 0.25
}

func (api *CexioApi) setBitcoin(amount float64) float64 {
	api.bitcoin += amount
	return api.bitcoin
}

func (api *CexioApi) setFunds(amount float64) float64 {
	api.funds += amount
	return api.funds
}
