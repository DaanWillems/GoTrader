package main

import (
	"fmt"
)

//Convert RAW OHLC Data to a neat map
func convertOHLC(cs []interface{}) map[string]string {
	return map[string]string{
		"time":  FloatToString(cs[0].(float64)), //Type asserts and then converts to string because its not officially a string
		"open":  cs[1].(string),
		"high":  cs[2].(string),
		"low":   cs[3].(string),
		"close": cs[4].(string),
	}
}

//Check all retrieved data, mainly intended for backtesting
func analyzeAll(data [][]interface{}) {
	fmt.Printf("\n Data: %v", data)
	for i := 0; i < 100; i++ {
		if i > len(data)-1 {
			fmt.Println("Exiting")
			return
		}
		checkHammer(data[i])
	}
}

//W.I.P, checks if a candlestick matches the hammer pattern.
func checkHammer(data []interface{}) {
	fmt.Printf("Open: %v ", convertOHLC(data)["open"])
	fmt.Printf("High: %v ", convertOHLC(data)["high"])
	fmt.Printf("Low: %v ", convertOHLC(data)["low"])
	fmt.Printf("Close: %v \n", convertOHLC(data)["close"])
}
