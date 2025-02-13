package main

import (
	"fmt"
	"main/src/exercise2/udp" // Defined msgTx(msg string) and msgRx()
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) > 1 {
		backup()
	} else {
		primary()
	}
}

func primary() {
	fmt.Println("--- Primary phase ---")
	// Start backup in new terminal window
	cmd := exec.Command("cmd.exe", "/C", "start", "cmd.exe", "/K", "go", "run", "main.go", "backup")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

	num := readNum()
	// Start sending heartbeats
	go func() {
		for {
			udp.SendMessage("127.0.0.1"+udp.BroadcastPort, "alive")
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Count 5 numbers
	for i := 1; i <= 5; i++ {
		num++
		writeNum(num)
		time.Sleep(time.Second)
	}

	// Exit current process
	os.Exit(0)
}

func backup() {
	fmt.Println("--- Backup phase ---")
	msgChan := udp.MsgRx()
	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Primary timeout - taking over")
			primary()
			return
		case msg := <-msgChan:
			if msg != "alive" {
				fmt.Println("Unexpected message - taking over")
				primary()
				return
			}
		}
	}
}

func writeNum(n int) {
	if os.WriteFile("testFile.txt", []byte(strconv.Itoa(n)), 0644) != nil {
		fmt.Println("Error writing to file")
	}
	fmt.Printf("Wrote number: %d\n", n)
}

func readNum() int {
	data, _ := os.ReadFile("testFile.txt")
	num, _ := strconv.Atoi(string(data))
	return num
}
