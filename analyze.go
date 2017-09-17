package main

import (
	"fmt"
	"math"
	"time"
)

var totalProfit float64 = 0
var position Position
var positionOpen = false
var tradeAmount float64 = 0.01

type Position struct {
	active    bool
	enteredAt time.Time
	leftAt    time.Time
	boughtFor float64
	soldFor   float64
}

//Compare two exchanges for arbitrage opportunities
func compare2(exchanges []Api, profit float64, count int) (float64, int) {
	var highest Api
	var lowest Api

	for _, e := range exchanges {
		e.RenewLast()
	}

	for i, e := range exchanges {
		if i == 0 {
			highest = e
			lowest = e
		}

		if e.GetLast() > highest.GetLast() {
			highest = e
		}

		if e.GetLast() < lowest.GetLast() {
			lowest = e
		}
	}

	fmt.Printf("Highest: %v %v, lowest: %v %v \n", highest.GetName(), highest.GetLast(), lowest.GetName(), lowest.GetLast())

	v1 := highest.GetLast()
	v2 := lowest.GetLast()

	var diff float64
	var currentProfit float64
	if v1 < v2 {
		diff = v2 - v1
		//fmt.Printf("%v is cheaper by %v %v %v \n", api1.GetName(), diff, v1, v2)
		currentProfit = diff * tradeAmount
		currentProfit = currentProfit - (currentProfit * highest.GetFee())
		currentProfit = currentProfit - (currentProfit * lowest.GetFee())
	} else if v2 < v1 {
		diff = v1 - v2

		currentProfit = diff * tradeAmount
		currentProfit = currentProfit - (currentProfit * lowest.GetFee())
		currentProfit = currentProfit - (currentProfit * highest.GetFee())
		//fmt.Printf("%v is cheaper by %v %v %v \n", api2.GetName(), diff, v1, v2)
	}

	if diff > 90 {
		fmt.Printf("Making trade!, trading %v bitcoin instantly results in %v euro's profit at %v with a difference of %v \n", tradeAmount, currentProfit, time.Now().Format(time.RFC850), diff)
		profit += currentProfit
		fmt.Printf("Total profit: %v \n", profit)
		fmt.Println("---------------------------------------------------------------------------------------------")
		highest.setBitcoin(tradeAmount - (tradeAmount * 2))
		highest.setFunds(profit)

		lowest.setBitcoin(tradeAmount)
		price := lowest.setBitcoin(0) * lowest.GetLast()
		lowest.setFunds(price - (price * 2))
		count++
	} else {
		fmt.Printf("Diff is %v, not trading \n", diff)
	}

	return profit, count
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

	if position.active {
		if (current["close"]-position.boughtFor > 30) || (current["close"]-position.boughtFor < -30) || (current["open"]-current["close"] > 0 && current["open"]-current["close"] < 15) {
			fmt.Println("Leaving position")
			position.leftAt = getTime(current["time"])
			position.soldFor = current["close"]

			diff := position.soldFor - position.boughtFor
			totalProfit += diff
			//Pretty print
			fmt.Printf("Bought at: %v\n", position.enteredAt)
			fmt.Printf("For: %v\n", position.boughtFor)
			fmt.Printf("Sold at: %v\n", position.leftAt)
			fmt.Printf("For: %v\n", position.soldFor)
			fmt.Printf("Difference is: %v\n", diff)
			fmt.Println("-------------------------------------------------------")
			position.active = false
			return
		}
	}

	//Check if current is a hammer
	var hammer = false
	var change float64
	var bottomWick float64
	var topWick float64
	//Get size
	if current["close"] > current["open"] {
		change = current["close"] - current["open"]
		bottomWick = current["open"] - current["low"]
		topWick = current["high"] - current["close"]
	} else {
		change = current["open"] - current["close"]
		bottomWick = current["close"] - current["low"]
		topWick = current["high"] - current["open"]
	}

	//Check regular hammer
	if current["close"] > current["open"] && (topWick/change)*100 < 20 && (bottomWick/change)*100 > 100 {
		hammer = true
	}

	// //Check inverted hammer
	// if current["close"] < current["open"] && (topWick/change)*100 > 110 && (bottomWick/change)*100 < 10 {
	// 	hammer = true
	// }

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

	fmt.Printf("Hammer approved at: %v %v %v %v %v %v \n", getTime(current["time"]), change, (topWick/change)*100, (bottomWick/change)*100, previousIsBullish, nextIsBullish)
	if position.active {
		if previousIsBullish && !nextIsBullish {
			fmt.Println("Leaving position")
			position.leftAt = getTime(current["time"])
			position.soldFor = current["close"]

			diff := position.soldFor - position.boughtFor
			totalProfit += diff
			//Pretty print
			fmt.Printf("Bought at: %v\n", position.enteredAt)
			fmt.Printf("For: %v\n", position.boughtFor)
			fmt.Printf("Sold at: %v\n", position.leftAt)
			fmt.Printf("For: %v\n", position.soldFor)
			fmt.Printf("Difference is: %v\n", diff)
			fmt.Println("-------------------------------------------------------")
			position.active = false
		}
		return
	}

	if previousIsBullish == false && nextIsBullish == true {
		if !(position.active) {
			fmt.Println("Establishing position")
			//If we arrive here everything looks good, so buy
			position = Position{
				active:    true,
				enteredAt: getTime(next["time"]),
				boughtFor: next["close"],
			}
		}
		return
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
