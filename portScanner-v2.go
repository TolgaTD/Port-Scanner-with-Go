package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var commonServices = map[int]string{
	20:    "FTP Data",
	21:    "FTP Control",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	67:    "DHCP Server",
	68:    "DHCP Client",
	69:    "TFTP",
	80:    "HTTP",
	110:   "POP3",
	119:   "NNTP",
	123:   "NTP",
	135:   "RPC",
	137:   "NetBIOS Name Service",
	138:   "NetBIOS Datagram Service",
	139:   "NetBIOS Session Service",
	143:   "IMAP",
	161:   "SNMP",
	179:   "BGP",
	389:   "LDAP",
	443:   "HTTPS",
	445:   "SMB",
	465:   "SMTPS",
	514:   "Syslog",
	543:   "Kerberos",
	544:   "Kerberos",
	548:   "AFP (Apple File Protocol)",
	554:   "RTSP",
	587:   "SMTP over TLS/SSL",
	631:   "Internet Printing Protocol",
	993:   "IMAPS",
	995:   "POP3S",
	1025:  "Microsoft RPC",
	1026:  "Windows RPC",
	1194:  "OpenVPN",
	1433:  "Microsoft SQL Server",
	1723:  "PPTP",
	3306:  "MySQL",
	3389:  "Remote Desktop Protocol",
	5060:  "SIP",
	5900:  "VNC",
	6001:  "X11",
	8000:  "Commonly used for HTTP",
	8080:  "HTTP Proxy",
	8443:  "HTTPS Proxy",
	8888:  "Commonly used for HTTP",
	10000: "Network Data Management Protocol",
}

func main() {
	clearConsole()

	// ASCII sanatını oluştur
	asciiArt := `
	_______ __   __ _______  ______  _____  __   _ _____ ______ _______       
	|______   \_/   |       |_____/ |     | | \  |   |    ____/ |______ |     
	______|    |    |_____  |    \_ |_____| |  \_| __|__ /_____ |______ |_____
																			  
   `

	// ASCII sanatını yazdır
	fmt.Println(asciiArt, "\n\n")

	// Başlangıç animasyonu
	animation := []string{"\\", "|", "/", "-"}
	for i := 0; i < 10; i++ {
		for _, frame := range animation {
			fmt.Printf("\rFast Port Scanner %s", frame)
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("\n\n")
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the target address to scan: ")
	scanner.Scan()
	targetAddress := scanner.Text()

	saveResult := ""
	for {
		fmt.Print("Do you want to save results? (YES: Y, NO: N): ")
		scanner.Scan()
		saveResult = strings.ToLower(scanner.Text())
		if saveResult == "y" || saveResult == "n" {
			break
		}
		fmt.Println("Invalid input, please enter 'Y' for Yes or 'N' for No.")
	}

	fmt.Println(`How do you want to scan?
1- Scan a range of ports (e.g., 0-100)
2- Scan specific ports (e.g., 20,23,80,443)
3- Scan the first 100 ports
4- Scan the first 1000 ports
5- Scan all ports`)
	fmt.Print("Enter your choice: ")
	scanner.Scan()
	scanOption := scanner.Text()

	startTime := time.Now()
	var startPort, endPort, openPorts, closedPorts int
	var specificPorts []int
	var file *os.File
	var err error

	if saveResult == "y" {
		formattedFileName := fmt.Sprintf("%s_scan_results.txt", targetAddress)
		file, err = os.Create(formattedFileName)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(file, "Scan Report for %s\n", targetAddress)
		fmt.Fprintf(file, "Scan started at: %s\n\n", startTime.Format(time.RFC1123))
	}

	switch scanOption {
	case "1", "2", "3", "4", "5":
		validateAndAssignPorts(scanner, &startPort, &endPort, &specificPorts, scanOption)
	default:
		fmt.Println("Invalid option, terminating the program.")
		return
	}

	limiter := rate.NewLimiter(100, 10) // Rate limiter: 10 requests per second

	setupScan(&wg, startPort, endPort, specificPorts, targetAddress, file, saveResult, limiter, &openPorts, &closedPorts)

	wg.Wait()

	duration := time.Since(startTime)
	fmt.Println("Scanning complete.")
	if saveResult == "y" {
		fmt.Fprintf(file, "\nScan completed in: %s\n", duration)
		fmt.Fprintf(file, "Total open ports: %d\n", openPorts)
		fmt.Fprintf(file, "Total closed/unreachable ports: %d\n", closedPorts)
	}
}
func promptForRange(scanner *bufio.Scanner) (int, int) {
	var startPort, endPort int
	var err error
	valid := false
	for !valid {
		fmt.Print("Enter the range of ports (e.g., 1-100): ")
		scanner.Scan()
		ports := strings.Split(scanner.Text(), "-")
		if len(ports) != 2 {
			fmt.Println("Invalid range format. Please use the format start-end (e.g., 1-100).")
			continue
		}
		startPort, err = strconv.Atoi(ports[0])
		endPort, err = strconv.Atoi(ports[1])
		if err != nil || startPort < 1 || endPort > 65535 || startPort > endPort {
			fmt.Println("Invalid port numbers. Ports must be between 1 and 65535 and start must be less than end.")
			continue
		}
		valid = true
	}
	return startPort, endPort
}
func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func promptForSpecificPorts(scanner *bufio.Scanner) []int {
	var ports []int
	valid := false
	for !valid {
		fmt.Print("Enter specific ports (e.g., 20,23,80,443): ")
		scanner.Scan()
		portStrs := strings.Split(scanner.Text(), ",")
		valid = true
		ports = []int{} // clear previous inputs
		for _, p := range portStrs {
			port, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil || port < 1 || port > 65535 {
				fmt.Printf("Invalid port number: %s. Ports must be between 1 and 65535.\n", p)
				valid = false
				break
			}
			ports = append(ports, port)
		}
	}
	return ports
}

func validateAndAssignPorts(scanner *bufio.Scanner, startPort, endPort *int, specificPorts *[]int, scanOption string) {
	switch scanOption {
	case "1":
		*startPort, *endPort = promptForRange(scanner)
	case "2":
		*specificPorts = promptForSpecificPorts(scanner)
	case "3":
		*startPort, *endPort = 1, 100
	case "4":
		*startPort, *endPort = 1, 1000
	case "5":
		*startPort, *endPort = 1, 65535
	}
}

func setupScan(wg *sync.WaitGroup, startPort, endPort int, specificPorts []int, targetAddress string, file *os.File, saveResult string, limiter *rate.Limiter, openPorts, closedPorts *int) {
	if len(specificPorts) > 0 {
		for _, port := range specificPorts {
			wg.Add(1)
			go scanPortWithRateLimit(port, targetAddress, file, wg, saveResult, limiter, openPorts, closedPorts)
		}
	} else {
		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			go scanPortWithRateLimit(port, targetAddress, file, wg, saveResult, limiter, openPorts, closedPorts)
		}
	}
}

func scanPortWithRateLimit(port int, targetAddress string, file *os.File, wg *sync.WaitGroup, saveResult string, limiter *rate.Limiter, openPorts, closedPorts *int) {
	defer wg.Done()
	limiter.Wait(context.Background()) // Wait to proceed within the rate limit
	address := fmt.Sprintf("%s:%d", targetAddress, port)
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		fmt.Printf("Port %d closed or unreachable: %v\n", port, err)
		if saveResult == "y" {
			fmt.Fprintf(file, "Port %d: Closed or Unreachable\n", port)
		}
		*closedPorts++
		return
	}
	defer conn.Close()
	if service, exists := commonServices[port]; exists {
		fmt.Printf("Port %d is open (Service: %s)\n\n\n\n", port, service)
		if saveResult == "y" {
			fmt.Fprintf(file, "Port %d: Open (Service: %s)\n\n\n\n", port, service)
		}
	} else {
		fmt.Printf("Port %d is open\n", port)
		if saveResult == "y" {
			fmt.Fprintf(file, "Port %d: Open\n", port)
		}
	}
	*openPorts++
}
