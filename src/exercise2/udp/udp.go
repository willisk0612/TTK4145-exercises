package udp

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"syscall"
	"context"
)

const (
	broadcastPort = ":30000"
	messagePort   = ":20009"
	broadcastIP   = "255.255.255.255"
	readTimeout   = 30 * time.Second
)

// Exported for use by other packages.
var BroadcastPort = broadcastPort

func handleError(err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatal(err)
		}
		log.Println(err)
	}
}

// RunUDP is used to launch continuous UDP messaging.
func RunUDP() {
	host, err := getOwnIP()
	handleError(err, false)

	sendIPBroadcast(host, broadcastIP+broadcastPort)
	msg := "Hello"
	go MsgTx(msg)
	go MsgRx()

	select {}
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

// SendMessage is the exported version of sendMessage.
func SendMessage(serverAddress, message string) {
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	handleError(err, true)

	conn, err := net.DialUDP("udp", nil, addr)
	handleError(err, true)
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	handleError(err, false)
}

func MsgTx(msg string) {
	for {
		time.Sleep(5 * time.Second)
		SendMessage("127.0.0.1"+broadcastPort, msg)
	}
}

func MsgRx() chan string {
	msgChan := make(chan string)
	go func() {
		// Use ListenConfig with SO_REUSEADDR enabled.
		lc := net.ListenConfig{
			Control: func(network, address string, c syscall.RawConn) error {
				var err error
				c.Control(func(fd uintptr) {
					err = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				})
				return err
			},
		}

		pc, err := lc.ListenPacket(context.Background(), "udp", broadcastPort)
		handleError(err, true)
		conn := pc.(*net.UDPConn)
		defer conn.Close()

		buffer := make([]byte, 1024)
		for {
			conn.SetReadDeadline(time.Now().Add(time.Second * 30))
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}
				handleError(err, true)
			}
			msgChan <- string(buffer[:n])
		}
	}()
	return msgChan
}
