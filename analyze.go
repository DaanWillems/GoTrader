package main

import (
	"fmt"
	"math"
	"time"
)

var totalProfit float64 = 0
var position Position
var positionOpen = false

type Position struct {
	active    bool
	enteredAt time.Time
	leftAt    time.Time
	boughtFor float64
	soldFor   float64
}

//Convert RAW OHLC Data to a neat map
func convertOHLC(cs []interface{}) map[string]float64 {
	return map[string]float64{
		"time":  cs[0].(float64), //Type asserts and then converts to string because its not officially a string
		"open":  StringtoFloat(cs[1].(string)),
		"high":  StringtoFloat(cs[2].(string)),
		"low":   StringtoFloat(cs[3].(string)),
		"close": StringtoFloat(cs[4].(string)),
	}
}

//Check all retrieved data, mainly intended for backtesting
func analyzeAll(data [][]interface{}) {

	position = Position{active: false}

	for i := 0; i < len(data)-1; i++ {
		Simulate(data, i)
	}
	fmt.Printf("Total profit over 1 months using 1 bitcoin with a 4 hour sell period: %v", totalProfit)
}

func Simulate(data [][]interface{}, i int) {

	//Get data
	current := convertOHLC(data[i])

	//Get previous
	if i-1 < 0 {
		return
	}
	previous := convertOHLC(data[i-1])

	//Get next
	if i+1 > len(data) {
		return
	}
	next := convertOHLC(data[i+1])

	//Check if current is a hammer
	var hammer = false

	//Check regular hammer
	if current["close"] > current["open"] && current["low"] < current["close"]-10 && current["high"] < current["open"]+20 {
		hammer = true
	}

	//Check inverted hammer
	if current["close"] > current["open"] && current["high"] > current["open"]+20 && current["low"] > current["close"]+10 {
		hammer = true
	}

	var nextIsBullish = false
	var previousIsBullish = false

	//Exit if there is no hammer
	if !hammer {
		return
	}

	//Check if previous was bearish
	if previous["close"] > previous["open"] {
		previousIsBullish = true
	}

	if next["close"] > next["open"] {
		nextIsBullish = true
	}

	if previousIsBullish && !nextIsBullish {
		if position.active {
			fmt.Println("Leaving position")
			diff := position.soldFor - position.boughtFor
			totalProfit += diff
			//Pretty print results
			fmt.Printf("Bought at: %v\n", position.enteredAt)
			fmt.Printf("For: %v\n", position.boughtFor)
			fmt.Printf("Sold at: %v\n", position.leftAt)
			fmt.Printf("For: %v\n", position.soldFor)
			fmt.Printf("Difference is: %v\n", diff)
			fmt.Println("-------------------------------------------------------")
			position.active = false
		}
	}

	if !previousIsBullish && nextIsBullish {
		if !(position.active) {
			fmt.Println("Establishing position")
			//If we arrive here everything looks good, so buy
			position = Position{
				active:    true,
				enteredAt: getTime(current["time"]),
				leftAt:    getTime(next["time"]),
				boughtFor: current["close"],
				soldFor:   next["close"],
			}
		}
	}

}

//W.I.P, checks if a candlestick matches the hammer pattern.
func checkHammer(data [][]interface{}, i int) bool {
	d := convertOHLC(data[i])
	// fmt.Printf("Open: %v ", convertOHLC(data)["open"])
	// fmt.Printf("High: %v ", convertOHLC(data)["high"])
	// fmt.Printf("Low: %v ", convertOHLC(data)["low"])
	// fmt.Printf("Close: %v \n", convertOHLC(data)["close"])

	sec, dec := math.Modf(d["time"])
	t := time.Unix(int64(sec), int64(dec*(1e9)))
	if d["close"] > d["open"] && d["low"] < d["close"]-5 {
		close := d["close"]
		buyprice := d["open"]
		if i-1 < 0 {
			return false
			fmt.Print("\n")
		}
		d = convertOHLC(data[i-1])
		if d["close"] < d["open"] && close < d["high"]+10 {
			close = d["close"]
			if i+1 > len(data) {
				return false
				fmt.Print("\n")
			}
			d = convertOHLC(data[i+1])
			diff := d["close"] - buyprice

			sec, dec := math.Modf(d["time"])
			t = time.Unix(int64(sec), int64(dec*(1e9)))
			fmt.Printf("Bought at: %v ", t)
			fmt.Printf("for: %v ", buyprice)
			fmt.Printf("Sold at: %v, for: %v ", t, d["close"])
			fmt.Printf("Profit: %v \n", diff)
			totalProfit = totalProfit + diff
			return true
		}
	}
	return false
}
