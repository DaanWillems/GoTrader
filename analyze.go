package main

import (
	"time"
)

type position struct {
	buyTime time.Time
	bitcoinAmount float64 //amount of bitcoin
	buyPrice float64
	sellBelow float64
	sellAbove float64 //Sell when the price reaches this target 
	
}

//Detects a hammer and returns false or true whether its advised to buy
func detectHammer() {

}
