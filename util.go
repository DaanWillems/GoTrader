package main

import (
	"math"
	"strconv"
	"time"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func StringtoFloat(input_string string) float64 {
	result, _ := strconv.ParseFloat(input_string[:len(input_string)-1], 64)
	return result
}

func getTime(input float64) time.Time {
	//Get time
	sec, dec := math.Modf(input)
	t := time.Unix(int64(sec), int64(dec*(1e9)))

	return t
}
