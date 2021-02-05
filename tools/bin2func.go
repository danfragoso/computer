package main

import (
	"io/ioutil"
	"fmt"
	"os"
)

func main() {
	payload, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
        fmt.Print(err)
		os.Exit(1)
	}

	fmt.Print("return []byte {")
	for _, b := range payload {
		fmt.Printf("0x%X, ", b)
	}
	fmt.Print("}")
}