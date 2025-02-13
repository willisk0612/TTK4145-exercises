package tcp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const (
	UDPBroadcastPort = 30000
	TCPFixedPort     = 34933
	TCPZeroPort      = 33546
	TCPListenPort    = 34934
	ReadTimeout      = 10 * time.Second
	ConnectDelay     = 100 * time.Millisecond
)

type Server struct {
	IP       string
	Listener net.Listener
}

func handleError(err error, fatal bool, msg string) {
	if err != nil {
		if fatal {
			log.Fatalf("%s: %v", msg, err)
		}
		log.Printf("%s: %v", msg, err)
	}
}

func RunTCP() {
	server, err := setupServer()
	handleError(err, true, "Server setup failed")
	defer server.Listener.Close()

	connections, err := connectToServers(server.IP)
	handleError(err, true, "Connection setup failed")
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	handleError(sendInitialMessages(connections[1]), false, "Failed to send initial messages")

	for {
		conn, err := server.Listener.Accept()
		handleError(err, false, "Accept error")
		if err == nil {
			go handleConnection(conn)
		}
	}
}

func setupServer() (*Server, error) {
	serverIP := make(chan string, 1)
	go listenForBroadcast(serverIP)

	ip := <-serverIP
	ownIP, err := getOwnIP()
	if err != nil {
		return nil, fmt.Errorf("IP lookup failed: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ownIP, TCPListenPort))
	if err != nil {
		return nil, fmt.Errorf("listener setup failed: %v", err)
	}

	return &Server{IP: ip, Listener: listener}, nil
}

func connectToServers(serverIP string) ([]net.Conn, error) {
	connections := make([]net.Conn, 0, 2)

	if conn1, err := connectToServer(serverIP, TCPFixedPort); err == nil {
		go handleConnection(conn1)
		connections = append(connections, conn1)
	}

	conn2, err := connectToServer(serverIP, TCPZeroPort)
	if err != nil {
		return nil, fmt.Errorf("zero-terminated connection failed: %v", err)
	}
	go handleConnection(conn2)
	connections = append(connections, conn2)

	return connections, nil
}

func sendInitialMessages(conn net.Conn) error {
	time.Sleep(ConnectDelay)
	ownIP, err := getOwnIP()
	if err != nil {
		return fmt.Errorf("IP lookup failed: %v", err)
	}

	messages := []string{
		fmt.Sprintf("Connect to: %s:%d\x00", ownIP, TCPListenPort),
		"Test echo\x00",
		"Hello server\x00",
	}

	for _, msg := range messages {
		if _, err := conn.Write([]byte(msg)); err != nil {
			return fmt.Errorf("failed to send %q: %v", msg, err)
		}
		time.Sleep(ConnectDelay)
	}
	return nil
}

func listenForBroadcast(serverIP chan<- string) {
	addr := fmt.Sprintf(":%d", UDPBroadcastPort)
	conn, err := net.ListenPacket("udp", addr)
	handleError(err, true, "UDP listen failed")
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buffer)
	handleError(err, true, "UDP read failed")

	msg := string(buffer[:n])
	parts := strings.Split(msg, " at ")
	if len(parts) != 2 {
		handleError(fmt.Errorf("invalid format"), true, "UDP message parsing failed")
	}

	ip := strings.TrimSuffix(parts[1], "!")
	fmt.Printf("Server IP: %s\n", ip)
	serverIP <- ip
}

func connectToServer(ip string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	fmt.Printf("Attempting to connect to: %s\n", addr)
	return net.Dial("tcp", addr)
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\x00')
		if err != nil {
			if err != io.EOF {
				handleError(err, false, "Read error")
			}
			return
		}
		trimmed := strings.TrimSpace(strings.TrimSuffix(msg, "\x00"))
		if trimmed != "" {
			fmt.Printf("Message received: %s\n", trimmed)
		}
	}
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
	return "", fmt.Errorf("no IP address found")
}
