// By @InfectedMunshroom
package linux

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
)

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

			if ip != nil && ip.IsLoopback() == false && ip.To4() != nil {
				return ip.String(), iface.HardwareAddr.String(), nil
			}
		}
	}

	return "", "", fmt.Errorf("no valid IP address found")
}

func ExtractPingStats(url string) (string, string, string, string, string, error) {
	cmd := exec.Command("ping", "-4", "-c", "5", url)
	output, err := cmd.Output()
	pingOutput := string(output)

	ipRe := regexp.MustCompile(`\((\d+\.\d+\.\d+\.\d+)\)`)
	statsRe := regexp.MustCompile(`rtt min/avg/max/mdev = ([\d.]+)/([\d.]+)/([\d.]+)/[\d.]+ ms`)
	packetLossRe := regexp.MustCompile(`(\d+)% packet loss`)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("Error in running program", err)
	}
	ipMatch := ipRe.FindStringSubmatch(pingOutput)
	if len(ipMatch) < 2 {
		return "", "", "", "", "", fmt.Errorf("no IP Address found")
	}
	ipAddress := ipMatch[1]

	statsMatch := statsRe.FindStringSubmatch(pingOutput)
	if len(statsMatch) < 4 {
		return "", "", "", "", "", fmt.Errorf("no timing statistics found")
	}
	minTime := statsMatch[1]
	avgTime := statsMatch[2]
	maxTime := statsMatch[3]

	packetLossMatch := packetLossRe.FindStringSubmatch(pingOutput)
	if len(packetLossMatch) < 2 {
		return "", "", "", "", "", fmt.Errorf("no packet loss information found")
	}
	packetLoss := packetLossMatch[1]

	return ipAddress, minTime, avgTime, maxTime, packetLoss, nil
}
