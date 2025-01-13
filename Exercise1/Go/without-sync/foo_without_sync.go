// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0

func incrementing() {
	//TODO: increment i 1000000 times
	for counter:=0;counter<1000000;counter++ {
		i++;
	}
}

func decrementing() {
	//TODO: decrement i 1000000 times
	for counter:=0;counter<1000000;counter++ {
		i--;
	}
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(2) // A: it makes 1 CPU run simontaneously instead

	// TODO: Spawn both functions as goroutines

	// Thread 1
	go incrementing()
	// Thread 2
	go decrementing()
	
	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	time.Sleep(500*time.Millisecond)
	Println("The magic number is:", i)
}
