package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// Get the server's own IP address
	host, err := getOwnIP()
	if err != nil {
		log.Fatalf("Error getting own IP address: %v", err)
	}

	// Broadcast the IP address on port 30000
	broadcastAddress := "255.255.255.255:34933"
	sendIPBroadcast(host, broadcastAddress)

	// errr := sendMessage("255.255.255.255:20009", "hello guys this is yohan")
	// if errr != nil {
	// 	log.Fatalf("Error sending message: %v", err)
	// }

	// // Start listening on port 30000 for messages
	// listenForMessages(":30000")
}

// Function to get the server's own local IP address
func getOwnIP() (string, error) {
	// Get a list of all network interfaces
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	// Iterate through interfaces and return the first non-loopback IP address
	for _, addr := range addrs {
		// Check for an IPv4 address
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("could not find an IP address")
}

// Function to broadcast the IP address on the specified address
func sendIPBroadcast(ipAddress, broadcastAddress string) {
	conn, err := net.Dial("tcp", broadcastAddress)
	if err != nil {
		log.Fatalf("Error dialing TCP: %v", err)
	}
	defer conn.Close()

	message := fmt.Sprintf("Server IP: %s", ipAddress)
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error sending broadcast: %v", err)
	}

	fmt.Printf("Broadcasted own IP address: %s\n", ipAddress)
}

func sendMessage(serverAddress string, message string) error {
	// Resolve the server address
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		return fmt.Errorf("error resolving server address: %v", err)
	}

	// Create a UDP connection to the server
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("error creating UDP connection: %v", err)
	}
	defer conn.Close()

	// Send the message to the server
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

// Function to listen for messages on the specified port
func listenForMessages(port string) {
	// Resolve the address to listen on
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}

	// Create a UDP connection to listen for incoming messages
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening on UDP port: %v", err)
	}
	defer conn.Close()

	// Set a read timeout for the UDP connection (optional)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	fmt.Println("Listening for incoming messages on port ", port, "...")

	// Listen for incoming messages
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if os.IsTimeout(err) {
				fmt.Println("Timeout reached, no more messages received.")
				break
			}
			log.Fatalf("Error reading from UDP: %v", err)
		}

		fmt.Printf("Received message: %s from %s\n", string(buffer[:n]), addr)
	}
}
