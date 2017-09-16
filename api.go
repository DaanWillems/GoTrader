package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

func newApi() *Api {

	api := &Api{}
	api.client = &http.Client{}

	api.pair = "XETHZEUR" //TODO: Hardcoded for now, should be moved to config file
	api.publicURL = "https://api.kraken.com/0/public"
	return api
}

type Api struct {
	publicURL string
	pair      string
	key       string
	secret    string
	client    *http.Client
}

type OHLC struct {
	Error  []interface{} `json:"error"`
	Result struct {
		Data [][]interface{} `json:"XETHZEUR"` //TODO: Hardcoded for now, should be moved to config file
		Last int             `json:"last"`
	} `json:"result"`
}

//Create request
func (api *Api) doRequest(parameters map[string]string, url string) []byte {
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
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
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

func (api *Api) getOHLC() OHLC {
	resp := api.doRequest(map[string]string{
		"interval": "60",
		"since":    "0",
	}, api.publicURL+"/OHLC")
	data := OHLC{}
	json.Unmarshal(resp, &data)
	return data
}
