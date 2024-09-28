package main

import (
	"fmt"
	"html/template"
	"net/http"
	"network-checker/internal/linux"
	"network-checker/internal/windows"
	"runtime"
)

type data struct {
	HostIP     string
	HostMAC    string
	PingIP     string
	MinTime    string
	AvgTime    string
	MaxTime    string
	PacketLoss string
}

func main() {
	// Set up the router
	http.HandleFunc("/", networkInfoHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func networkInfoHandler(w http.ResponseWriter, r *http.Request) {
	os := runtime.GOOS
	var output data

	if os == "linux" {
		output.HostIP, output.HostMAC, output.PingIP, output.MinTime, output.AvgTime, output.MaxTime, output.PacketLoss = getLinuxNetworkInfo()
	} else if os == "windows" {
		output.HostIP, output.HostMAC, output.PingIP, output.MinTime, output.AvgTime, output.MaxTime, output.PacketLoss = getWindowsNetworkInfo()
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, output)
}

func getLinuxNetworkInfo() (string, string, string, string, string, string, string) {
	hostIP, hostMAC, err := linux.GetHostIPAndMAC()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	ip, minTime, avgTime, maxTime, packetLoss, err := linux.ExtractPingStats("google.com")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// Convert packet loss to percentage
	packetLossPercentage := packetLoss + " %"
	return hostIP, hostMAC, ip, minTime, avgTime, maxTime, packetLossPercentage
}

func getWindowsNetworkInfo() (string, string, string, string, string, string, string) {
	hostIP, hostMAC, err := windows.GetHostIPAndMAC()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	ip, minTime, avgTime, maxTime, packetLoss, err := windows.ExtractPingStats("google.com")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// Convert packet loss to percentage
	packetLossPercentage := packetLoss + " %"
	return hostIP, hostMAC, ip, minTime, avgTime, maxTime, packetLossPercentage
}
