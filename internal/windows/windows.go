// By @sanyaw25
package windows

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

			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), iface.HardwareAddr.String(), nil
			}
		}
	}

	return "", "", fmt.Errorf("no valid IP address found")
}

func ExtractPingStats(url string) (string, string, string, string, string, error) {
	cmd := exec.Command("ping", "-4", "-n", "5", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error running ping command: %w", err)
	}

	pingOutput := string(output)
	fmt.Println("Ping Command Output:\n", pingOutput)

	ipRe := regexp.MustCompile(`Pinging ([\w.-]+) \[(\d+\.\d+\.\d+\.\d+)\]`)
	statsRe := regexp.MustCompile(`Minimum = (\d+)ms, Maximum = (\d+)ms, Average = (\d+)ms`)
	packetLossRe := regexp.MustCompile(`Lost = (\d+) \((\d+)% loss\)`)

	ipMatch := ipRe.FindStringSubmatch(pingOutput)
	if len(ipMatch) < 3 {
		return "", "", "", "", "", fmt.Errorf("no IP address found")
	}
	ipAddress := ipMatch[2]

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
