
package main

import (
 "bufio"
 "fmt"
 "net"
 "os"
) 

func startClient(address string) {
 //connect to this socket
 connClient, _ := net.Dial("tcp", address)

for {

    //read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')

    //send to socket
    fmt.Fprint(connClient, text+"\n")

    //listen for reply
    //message, _ := bufio.NewReader(connClient).ReadString('\n')
    //fmt.Print("Message from server: " + message)
  }
 }
func startServer() {
 fmt.Println("Starting...")

//listen on all interfaces
ln, _ := net.Listen("tcp", ":8081")

//accept connection on port
connServer, _ := ln.Accept()

//run loop forever
for {
    //will listen for message to process ending in newline(\n)
    message, _ := bufio.NewReader(connServer).ReadString('\n')

    connServer.Write([]byte(message + "\n"))
  }
 }