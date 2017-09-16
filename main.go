package main

import (
	"fmt"
)

func main() {
	api := newApi()
	resp := api.getOHLC()
	fmt.Print("Last: ")
	fmt.Println(resp.Result.Last)
	analyzeAll(resp.Result.Data)
}
