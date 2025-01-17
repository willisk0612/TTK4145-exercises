// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

type Operation int

const (
	Nothing Operation = iota // default
	Increment
	Decrement
)

type Request struct {
	op   Operation
	get  chan int
	done chan bool
}

func numberServer(requests chan Request) {
	i := 0
	for {
		select {
		case req := <-requests:
			switch req.op {
			case Increment:
				i++
				req.done <- true
			case Decrement:
				i--
				req.done <- true
			default:
				req.get <- i
			}
		}
	}
}

func incrementing(reqs chan Request, done chan bool) {
	//TODO: increment i 1000000 times
	for counter := 0; counter < 1000000; counter++ {
		reqs <- Request{op: Increment}
	}
	reqs <- Request{done: done}
}

func decrementing(reqs chan Request, done chan bool) {
	//TODO: decrement i 1000000 times
	for counter := 0; counter < 1000000; counter++ {
		reqs <- Request{op: Decrement}
	}
	reqs <- Request{done: done}
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(2) // A: It makes two goroutines run parallell. Setting it to 1 would make the program run sequentially.

	// TODO: Spawn both functions as goroutines
	reqs := make(chan Request)
	done := make(chan bool)

	go incrementing(reqs, done)
	go decrementing(reqs, done)
	go numberServer(reqs)

	<-done
	<-done

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	time.Sleep(500 * time.Millisecond)

	result := make(chan int)
	reqs <- Request{get: result}
	Println("The magic number is:", <-result)
}
