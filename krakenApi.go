package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func newKrakenApi() *KrakenApi {

	api := &KrakenApi{}
	api.client = &http.Client{}

	api.pair = "XXBTZEUR" //TODO: Hardcoded for now, should be moved to config file
	api.publicURL = "https://api.kraken.com/0/public"

	api.funds = 300
	api.bitcoin = 0.03

	return api
}

type KrakenApi struct {
	publicURL string
	pair      string
	key       string
	secret    string
	client    *http.Client
	last      float64
	bitcoin   float64
	funds     float64
}

type KrakenTickers struct {
	Error  []interface{} `json:"error"`
	Result struct {
		XXBTZEUR struct {
			A []string `json:"a"`
			B []string `json:"b"`
			C []string `json:"c"`
			V []string `json:"v"`
			P []string `json:"p"`
			T []int    `json:"t"`
			L []string `json:"l"`
			H []string `json:"h"`
			O string   `json:"o"`
		} `json:"XXBTZEUR"`
	} `json:"result"`
}

type KrakenOHLC struct {
	Error  []interface{} `json:"error"`
	Result struct {
		Data [][]interface{} `json:"XXBTZEUR"` //TODO: Hardcoded for now, should be moved to config file
		Last int             `json:"last"`
	} `json:"result"`
}

//Create request
func (api *KrakenApi) doRequest(parameters map[string]string, url string) []byte {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed")
		log.Print(err)
	}
	req.Header.Add("User-Agent", "GoTrade")
	q := req.URL.Query()

	q.Add("pair", api.pair)

	for key, value := range parameters {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()
	// Save a copy of this request for debugging.
	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(requestDump))
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

func (api *KrakenApi) getOHLC() KrakenOHLC {
	resp := api.doRequest(map[string]string{
		"interval": "21600",
		"since":    "0",
	}, api.publicURL+"/OHLC")
	data := KrakenOHLC{}
	json.Unmarshal(resp, &data)
	return data
}

func (api *KrakenApi) RenewLast() float64 {
	resp := api.doRequest(nil, api.publicURL+"/Ticker")
	data := KrakenTickers{}
	json.Unmarshal(resp, &data)
	api.last = StringtoFloat(data.Result.XXBTZEUR.C[0])
	return StringtoFloat(data.Result.XXBTZEUR.C[0])
}

func (api *KrakenApi) GetLast() float64 {
	return api.last
}

func (api *KrakenApi) GetName() string {
	return "Kraken"
}

func (api *KrakenApi) GetFee() float64 {
	return 0.16
}

func (api *KrakenApi) setBitcoin(amount float64) float64 {
	api.bitcoin += amount
	return api.bitcoin
}

func (api *KrakenApi) setFunds(amount float64) float64 {
	api.funds += amount
	return api.funds
}
