package windows

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
)

// GetHostIPAndMAC retrieves the first non-loopback IP address and MAC address of the system
func GetHostIPAndMAC() (string, string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", "", err
	}

	for _, iface := range interfaces {
		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Check if it's a valid IPv4 address, skipping loopback
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), iface.HardwareAddr.String(), nil
			}
		}
	}

	return "", "", fmt.Errorf("no valid IP address found")
}

// ExtractPingStats runs a ping to a given URL and extracts min/avg/max times and packet loss
func ExtractPingStats(url string) (string, string, string, string, string, error) {
	// Windows uses -n for the count of ping requests
	cmd := exec.Command("ping", "-n", "5", url)
	output, err := cmd.CombinedOutput() // Use CombinedOutput to get stderr and stdout
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error running ping command: %w", err)
	}

	pingOutput := string(output)
	fmt.Println("Ping Command Output:\n", pingOutput) // Display raw output

	// Regular expression for extracting the IP address from the ping output
	ipRe := regexp.MustCompile(`Pinging ([\w.-]+) \[(\d+\.\d+\.\d+\.\d+)\]`)
	// Regular expression for extracting the min/avg/max ping times
	statsRe := regexp.MustCompile(`Minimum = (\d+)ms, Maximum = (\d+)ms, Average = (\d+)ms`)
	// Regular expression for extracting packet loss information
	packetLossRe := regexp.MustCompile(`Lost = (\d+) \((\d+)% loss\)`)

	ipMatch := ipRe.FindStringSubmatch(pingOutput)
	if len(ipMatch) < 3 {
		return "", "", "", "", "", fmt.Errorf("no IP address found")
	}
	ipAddress := ipMatch[2] // The IP address is in the second capturing group

	statsMatch := statsRe.FindStringSubmatch(pingOutput)
	if len(statsMatch) < 4 {
		return "", "", "", "", "", fmt.Errorf("no timing statistics found")
	}
	minTime := statsMatch[1]
	maxTime := statsMatch[2]
	avgTime := statsMatch[3]

	packetLossMatch := packetLossRe.FindStringSubmatch(pingOutput)
	if len(packetLossMatch) < 3 {
		return "", "", "", "", "", fmt.Errorf("no packet loss information found")
	}
	packetLoss := packetLossMatch[2]

	return ipAddress, minTime, avgTime, maxTime, packetLoss, nil
}

// main function to test the functionalities
// func main() {
// ip, mac, err := GetHostIPAndMAC()
// if err != nil {
// fmt.Println("Error getting IP and MAC:", err)
// return
// }
// fmt.Println("IP Address:", ip)
// fmt.Println("MAC Address:", mac)
//
// Change the URL to a valid address for pinging
// url := "google.com" // Example URL
// ipAddress, minTime, avgTime, maxTime, packetLoss, err := ExtractPingStats(url)
// if err != nil {
// fmt.Println("Error extracting ping stats:", err)
// return
// }
//
// fmt.Printf("Ping Statistics for %s:\n", url)
// fmt.Printf("IP Address: %s\n", ipAddress)
// fmt.Printf("Minimum Time: %s ms\n", minTime)
// fmt.Printf("Average Time: %s ms\n", avgTime)
// fmt.Printf("Maximum Time: %s ms\n", maxTime)
// fmt.Printf("Packet Loss: %s%%\n", packetLoss)
// }
//
