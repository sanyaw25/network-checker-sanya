package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type data struct {
	mac               string
	hostname          string
	hostip            string
	urlcust           string
	urlgoogle         string
	packetLossPercent int
	timeTakenMin      float32
	timeTakenAvg      float32
	timeTakenMax      float32
}

// Function to execute the ping command and return the output
func pingURL(url string) (string, error) {
	cmd := exec.Command("ping", "-c", "5", url)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func main() {
	// Prompt the user for input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the website/URL to ping: ")
	url, _ := reader.ReadString('\n')
	url = url[:len(url)-1] // Remove the newline character

	// Call the pingURL function
	result, err := pingURL(url)
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	// Print the output
	fmt.Println("Command Output:")
	fmt.Println(result)
}
