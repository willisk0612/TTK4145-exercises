package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	broadcastPort = ":30000"
	messagePort   = ":20009"
	broadcastIP   = "255.255.255.255"
	readTimeout   = 30 * time.Second
)

func handleError(err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatal(err)
		}
		log.Println(err)
	}
}

func main() {
	host, err := getOwnIP()
	handleError(err, false)

	sendIPBroadcast(host, broadcastIP+broadcastPort)
	listenForMessages(broadcastPort)
}

func getOwnIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("no suitable IP address found")
}

func sendIPBroadcast(ipAddress, broadcastAddress string) {
	conn, err := net.Dial("udp", broadcastAddress)
	handleError(err, true)
	defer conn.Close()

	message := fmt.Sprintf("Server IP: %s", ipAddress)
	_, err = conn.Write([]byte(message))
	handleError(err, true)

	log.Printf("Broadcasted own IP address: %s\n", ipAddress)
}

func sendMessage(serverAddress, message string) {
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	handleError(err, true)

	conn, err := net.DialUDP("udp", nil, addr)
	handleError(err, true)
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	handleError(err, false)
}

func listenForMessages(port string) {
	addr, err := net.ResolveUDPAddr("udp", port)
	handleError(err, true)

	conn, err := net.ListenUDP("udp", addr)
	handleError(err, true)
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(readTimeout))
	log.Printf("Listening for incoming messages on port %s...\n", port)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if os.IsTimeout(err) {
				log.Println("Timeout reached, no more messages received.")
				return
			}
			handleError(err, true)
		}
		log.Printf("Received message: %s from %s\n", string(buffer[:n]), addr)
	}
}
