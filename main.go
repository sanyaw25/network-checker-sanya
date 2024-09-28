package main

import (
	"fmt"
	"network-checker/internal/linux"
	"runtime"
)

type data struct {
	hostIP     string
	hostMAC    string
	ip         string
	minTime    string
	avgTime    string
	maxTime    string
	packetLoss string
	err        error
}

func main() {
	os := runtime.GOOS

	fmt.Println("OS:", os)

	var output data

	if os == "linux" {
		fmt.Println("Using Linux based script")
		output.hostIP, output.hostMAC, output.err = linux.GetHostIPAndMAC()

		if output.err != nil {
			fmt.Println("Error: ", output.err)
		}
		output.ip, output.minTime, output.avgTime, output.maxTime, output.packetLoss, output.err = linux.ExtractPingStats("google.com")
		if output.err != nil {
			fmt.Println("Error: ", output.err)
			return
		}

		//Statements to check stuff is getting recieved
		fmt.Printf("Host IP: %s\nHost MAC: %s\n", output.hostIP, output.hostMAC)
		fmt.Printf("IP: %s\nMin Time: %s\nAvg Time: %s\n", output.ip, output.minTime, output.avgTime)

	}

}
