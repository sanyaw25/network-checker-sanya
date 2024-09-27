package main

import (
	"fmt"
	"network-checker/internal/linux"
)

func main() {
	fmt.Println("Hi, this is working")
	result, err := linux.PingURL("google.com")
	if err != nil {
		fmt.Println("Some error occured")
		return
	}

	fmt.Println(result)

}
