package main

import (
	"fmt"
)


func main() {
	api := newApi()
	resp := api.getOHLC()
	fmt.Printf("Result: %v", resp.Result.XXBTZEUR[1])
}

