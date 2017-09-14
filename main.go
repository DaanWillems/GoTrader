package main

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/beldur/kraken-go-api-client"
)

func main() {
	api := krakenapi.New("KEY", "SECRET")
	result, err := api.Query("OHLC", map[string]string{
		"pair": "XXBTZEUR",
		"since": "1",
	})

	if err != nil {
		log.Fatal(err)
	}

	json, _ := json.Marshal(result)

	fmt.Printf("Result: %+v\n", string(json))

	// There are also some strongly typed methods available
	ticker, err := api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ticker.XXBTZEUR.OpeningPrice)
}
