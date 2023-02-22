package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"os"
)

func main() {
	// Get the IP range and port list from the command line arguments
	ipRange := ""
	ports := ""

	if len(os.Args) >= 2 {
		ipRange = os.Args[1]
	}
	if len(os.Args) >= 3 {
		ports = os.Args[2]
	}

	// Split the IP range into its start and end IP addresses
	ipStart, ipEnd, err := parseIPRange(ipRange)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse the port list into individual ports
	portList, err := parsePortList(ports)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a wait group to ensure all goroutines finish before exiting
	var wg sync.WaitGroup

	// Loop over all IP addresses in the range and scan each one for open ports
	for ip := ipStart; ip <= ipEnd; ip++ {
		ipStr := intToIPString(ip)

		for _, port := range portList {
			wg.Add(1)
			go func(ip string, port int) {
				defer wg.Done()

				// Connect to the server
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 1e9)
				if err != nil {
					return
				}
				conn.Close()

				// Print the open port number
				fmt.Printf("Port %d is open on %s\n", port, ip)
			}(ipStr, port)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// parseIPRange parses an IP range string into start and end IP addresses
func parseIPRange(ipRange string) (uint32, uint32, error) {
	if ipRange == "" {
		return 0, 0, fmt.Errorf("IP range must be specified")
	}

	// Split the IP range into its start and end IP addresses
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Invalid IP range: %s", ipRange)
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	if startIP == nil {
		return 0, 0, fmt.Errorf("Invalid start IP address: %s", parts[0])
	}

	endIP := net.ParseIP(strings.TrimSpace(parts[1]))
	if endIP == nil {
		return 0, 0, fmt.Errorf("Invalid end IP address: %s", parts[1])
	}

	startIPInt := ipToInt(startIP)
	endIPInt := ipToInt(endIP)

	if startIPInt > endIPInt {
		return 0, 0, fmt.Errorf("Start IP address cannot be greater than end IP address")
	}

	return startIPInt, endIPInt, nil
}

// parsePortList parses a port list string into individual ports
func parsePortList(portList string) ([]int, error) {
	if portList == "" {
		return nil, fmt.Errorf("Port list must be specified")
	}

	// Split the port list into individual ports
	portParts := strings.Split(portList, ",")
	ports := []int{}

	for _, portPart := range portParts {
		portRange := strings.Split(portPart, "-")

// ...

		if len(portRange) == 1 {
			portNum, err := strconv.Atoi(portRange[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: %s", portRange[0])
			}

			if portNum < 1 || portNum > 65535 {
				return nil, fmt.Errorf("Port number out of range: %d", portNum)
			}

			ports = append(ports, portNum)
		} else if len(portRange) == 2 {
			startPortNum, err := strconv.Atoi(portRange[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid start port number: %s", portRange[0])
			}

			endPortNum, err := strconv.Atoi(portRange[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid end port number: %s", portRange[1])
			}

			if startPortNum < 1 || startPortNum > 65535 {
				return nil, fmt.Errorf("Start port number out of range: %d", startPortNum)
			}
			if endPortNum < 1 || endPortNum > 65535 {
				return nil, fmt.Errorf("End port number out of range: %d", endPortNum)
			}
			if startPortNum > endPortNum {
				return nil, fmt.Errorf("Start port number cannot be greater than end port number")
			}

			for i := startPortNum; i <= endPortNum; i++ {
				ports = append(ports, i)
			}
		} else {
			return nil, fmt.Errorf("Invalid port range: %s", portPart)
		}
	}

	return ports, nil
}

// intToIPString converts an integer representation of an IP address to a string
func intToIPString(ip uint32) string {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}

// ipToInt converts an IP address to its integer representation
func ipToInt(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
