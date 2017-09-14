package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var PublicURL = "api.kraken.com/0/public"

func newApi() *Api {

	api := &Api{}
	api.client = &http.Client{}

	api.pairs = "XXBTZEUR"

	return api
}

type Api struct {
	pairs  string
	key    string
	secret string
	client *http.Client
}

//Create request
func (api *Api) doRequest(parameters map[string]string, url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request failed")
		log.Print(err)
	}

	q := req.URL.Query()

	q.Add("pairs", api.pairs)

	for key, value := range parameters {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

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

func (api *Api) getOHLC() {

}
