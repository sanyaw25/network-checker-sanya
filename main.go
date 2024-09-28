package main

import (
	"fmt"
	"network-checker/internal/linux"
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
	//Statement to check if stuff is working
	fmt.Println("Hi, this is working")

	var output data

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
