package main

import (
)

func main() {
	api := []Api{
		newBitstampApi(),
		newKrakenApi(),
		newCexioApi(),
	}
	var profit float64
	var count int
	for {
		profit, count = compare2(api, profit, count)
		duration := time.Duration(5) * time.Second
		time.Sleep(duration)
		if count >= 10 {
			break
		}
	}

	for _, e := range api {
		fmt.Printf("Name %v, bitcoins: %v, funds: %v\n ", e.GetName(), e.setBitcoin(0), e.setFunds(0))
	}
	fmt.Println("Out of bitcoin!")
}
